package server

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/to77e/word-of-wisdom/internal/message"
	"github.com/to77e/word-of-wisdom/internal/proofofwork"
	"github.com/to77e/word-of-wisdom/internal/wordofwisdom"
	"github.com/to77e/word-of-wisdom/tools/validator"
)

const (
	numberOfWorkers   = 2
	defaultBufferSize = 128
	defaultTimeout    = time.Second
)

type Server struct {
	listener   net.Listener
	quit       chan struct{}
	wg         sync.WaitGroup
	connection chan net.Conn
}

func New(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("tcp listener: %s\n", err.Error())
	}
	slog.With("address", address).Info("tcp listen")

	return &Server{
		listener:   listener,
		quit:       make(chan struct{}),
		connection: make(chan net.Conn),
	}, nil
}

func (s *Server) Start() {
	s.wg.Add(numberOfWorkers)
	go s.acceptConnections()
	go s.handleConnections()
}

func (s *Server) Stop() {
	close(s.quit)
	if err := s.listener.Close(); err != nil {
		slog.With("error", err.Error()).Error("close tcp listener")
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("tcp listener closed")
		return
	case <-time.After(defaultTimeout):
		slog.Info("timed out waiting for connections to finish")
		return
	}
}

func (s *Server) acceptConnections() {
	// todo: gracefully stop accept connection
	defer s.wg.Done()
	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				slog.With("error", err.Error()).Error("accept tcp connection")
				return
			}
			s.connection <- conn
		}
	}
}

func (s *Server) handleConnections() {
	defer s.wg.Done()
	for {
		select {
		case <-s.quit:
			return
		case conn := <-s.connection:
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.With("error", err.Error()).Error("close tcp connection")
		}
		slog.With("remote_address", conn.RemoteAddr().String()).Info("tcp connection closed")
	}()

	// implementing the Proof of Work consensus mechanism
	// (https://en.wikipedia.org/wiki/Proof_of_work)
	// with the challenge-response protocol

	// 1. request service
	if err := s.ReadChallengeRequest(conn); err != nil {
		s.HandleError(conn, err, "read challenge request")
		return
	}

	// 2. choose
	pow := proofofwork.New()
	if err := pow.GenerateChallenge(); err != nil {
		s.HandleError(conn, err, "generate challenge")
		return
	}

	// 3. challenge
	if err := s.WriteChallengeResponse(conn, pow.GetChallenge()); err != nil {
		s.HandleError(conn, err, "write challenge response")
		return
	}

	// 5. response
	solution, err := s.ReadSolutionRequest(conn)
	if err != nil {
		s.HandleError(conn, err, "read solution request")
		return
	}

	// 6. verify
	pow.SetSolution(solution)
	isVerified, err := pow.VerifySolution()
	if err != nil {
		s.HandleError(conn, err, "verify solution")
		return
	}
	if !isVerified {
		s.HandleError(conn, err, "solution is not verified")
		return
	}

	solutionResponse := &SolutionMessageResponse{
		Type:   SolutionResponse,
		Result: wordofwisdom.GetRandomQuote(),
	}

	// 7. grant service
	if err = s.sendResponse(conn, solutionResponse); err != nil {
		s.HandleError(conn, err, "send response")
		return
	}
}

func (s *Server) ReadChallengeRequest(conn net.Conn) error {
	challengeRequest := &ChallengeMessageRequest{}
	if err := s.getRequest(conn, challengeRequest); err != nil {
		return fmt.Errorf("get request: %w\n", err)
	}
	return nil
}

func (s *Server) WriteChallengeResponse(conn net.Conn, challenge []byte) error {
	challengeResponse := &ChallengeMessageResponse{
		Type:      ChallengeResponse,
		Challenge: challenge,
	}
	if err := s.sendResponse(conn, challengeResponse); err != nil {
		return fmt.Errorf("send response: %w\n", err)
	}
	return nil
}

func (s *Server) ReadSolutionRequest(conn net.Conn) ([]byte, error) {
	solutionRequest := &SolutionMessageRequest{
		Solution: make([]byte, 0),
	}
	if err := s.getRequest(conn, solutionRequest); err != nil {
		return nil, fmt.Errorf("get request: %w\n", err)
	}
	return solutionRequest.Solution, nil
}

func (s *Server) HandleError(conn net.Conn, err error, message string) {
	errorMessage := &ErrorMessage{
		Type:         Error,
		ErrorMessage: message,
	}
	if errResp := s.sendResponse(conn, errorMessage); errResp != nil {
		slog.With("error", errResp.Error()).Error("send error response")
	}
	slog.With("error", err.Error()).Error(message)
}

func (s *Server) getRequest(conn net.Conn, req interface{}) error {
	// todo: dynamic buffer size
	buf := make([]byte, defaultBufferSize)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("read tcp connection: %w\n", err)
	}
	slog.With("bytes", n, "remote_address", conn.RemoteAddr().String()).Info("tcp connection")

	m := message.New(buf[:n:n], validator.Get())
	if err = m.UnmarshalData(req); err != nil {
		return fmt.Errorf("unmarshal data: %w\n", err)
	}

	return nil
}

func (s *Server) sendResponse(conn net.Conn, data interface{}) error {
	m := message.New(nil, nil)
	if err := m.MarshalData(data); err != nil {
		return fmt.Errorf("marshal data: %w\n", err)
	}
	n, err := conn.Write(m.GetData())
	if err != nil {
		return fmt.Errorf("write tcp connection: %w\n", err)
	}
	slog.With("bytes", n, "remote_address", conn.RemoteAddr().String()).Info("tcp connection wrote")
	return nil
}
