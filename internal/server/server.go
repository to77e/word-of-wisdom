package server

import (
	"errors"
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
	numberOfWorkers          = 2
	defaultBufferSize        = 128
	defaultTimeout           = time.Second
	defaultConnectionTimeout = time.Second * 10
	defaultDifficulty        = 1
)

type ProofOfWorker interface {
	GenerateChallenge() error
	GetChallenge() []byte
	SetSolution(solution []byte)
	GetDifficulty() int
	VerifySolution() bool
}

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
			if err = conn.SetDeadline(time.Now().Add(defaultConnectionTimeout)); err != nil {
				slog.With("error", err.Error()).Error("set deadline")
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
			pow := proofofwork.New(defaultDifficulty)
			go s.handleConnection(conn, pow)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn, pow ProofOfWorker) {
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
	if err := s.readChallengeRequest(conn); err != nil {
		s.handleError(conn, err, "read challenge request")
		return
	}

	// 2. choose
	if err := pow.GenerateChallenge(); err != nil {
		s.handleError(conn, err, "generate challenge")
		return
	}

	// 3. challenge
	if err := s.writeChallengeResponse(conn, pow.GetChallenge(), pow.GetDifficulty()); err != nil {
		s.handleError(conn, err, "write challenge response")
		return
	}

	// 5. response
	solution, err := s.readSolutionRequest(conn)
	if err != nil {
		s.handleError(conn, err, "read solution request")
		return
	}

	// 6. verify
	pow.SetSolution(solution)
	isVerified := pow.VerifySolution()
	if !isVerified {
		s.handleError(conn, errors.New("solution is not valid"), "verify solution")
		return
	}

	solutionResponse := &message.SolutionMessageResponse{
		Type:  message.SolutionResponse,
		Quote: wordofwisdom.GetRandomQuote(),
	}

	// 7. grant service
	if err = s.sendResponse(conn, solutionResponse); err != nil {
		s.handleError(conn, err, "send response")
		return
	}
}

func (s *Server) readChallengeRequest(conn net.Conn) error {
	challengeRequest := &message.ChallengeMessageRequest{}
	if err := s.getRequest(conn, challengeRequest); err != nil {
		return fmt.Errorf("get request: %w\n", err)
	}
	return nil
}

func (s *Server) writeChallengeResponse(conn net.Conn, challenge []byte, difficulty int) error {
	challengeResponse := &message.ChallengeMessageResponse{
		Type:       message.ChallengeResponse,
		Challenge:  challenge,
		Difficulty: difficulty,
	}
	if err := s.sendResponse(conn, challengeResponse); err != nil {
		return fmt.Errorf("send response: %w\n", err)
	}
	return nil
}

func (s *Server) readSolutionRequest(conn net.Conn) ([]byte, error) {
	solutionRequest := &message.SolutionMessageRequest{
		Solution: make([]byte, 0),
	}
	if err := s.getRequest(conn, solutionRequest); err != nil {
		return nil, fmt.Errorf("get request: %w\n", err)
	}
	return solutionRequest.Solution, nil
}

func (s *Server) handleError(conn net.Conn, err error, content string) {
	errorMessage := &message.ErrorMessage{
		Type:         message.Error,
		ErrorMessage: content,
	}
	if errResp := s.sendResponse(conn, errorMessage); errResp != nil {
		slog.With("error", errResp.Error()).Error("send error response")
	}
	slog.With("error", err.Error()).Error(content)
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
