package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/to77e/word-of-wisdom/internal/config"
	"github.com/to77e/word-of-wisdom/internal/messages"
	"github.com/to77e/word-of-wisdom/internal/proofofwork"
	"github.com/to77e/word-of-wisdom/internal/wordofwisdom"
	"github.com/to77e/word-of-wisdom/tools/validator"
)

const (
	numberOfWorkers = 2
	defaultTimeout  = time.Second
)

type ProofOfWorker interface {
	GenerateChallenge() error
	GetChallenge() []byte
	SetSolution(solution []byte)
	GetDifficulty() int
	VerifySolution() bool
}

type Server struct {
	listener    net.Listener
	quit        chan struct{}
	wg          sync.WaitGroup
	connections chan net.Conn
}

func New(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("tcp listener: %s\n", err.Error())
	}

	slog.With("address", address).Info("tcp listen")

	return &Server{
		listener:    listener,
		quit:        make(chan struct{}),
		connections: make(chan net.Conn),
	}, nil
}

func (s *Server) Start() {
	cfg := config.GetInstance()

	s.wg.Add(numberOfWorkers)
	go s.acceptConnections(cfg.Server)
	go s.handleConnections(cfg.ProofOfWork)
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

func (s *Server) acceptConnections(cfg config.Server) {
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
			if err = conn.SetDeadline(time.Now().Add(cfg.ConnectionTimeout)); err != nil {
				slog.With("error", err.Error()).Error("set deadline")
				return
			}
			s.connections <- conn
		}
	}
}

func (s *Server) handleConnections(cfg config.ProofOfWork) {
	defer s.wg.Done()
	for {
		select {
		case <-s.quit:
			return
		case conn := <-s.connections:
			go s.handleConnection(conn, proofofwork.New(cfg.Difficulty))
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
	if !pow.VerifySolution() {
		s.handleError(conn, errors.New("solution is not valid"), "verify solution")
		return
	}

	// 7. grant service
	solutionResponse := &messages.SolutionMessageResponse{
		Type:  messages.SolutionResponse,
		Quote: wordofwisdom.GetRandomQuote(),
	}
	if err = s.sendResponse(conn, solutionResponse); err != nil {
		s.handleError(conn, err, "send response")
		return
	}
}

func (s *Server) readChallengeRequest(conn net.Conn) error {
	challengeRequest := &messages.ChallengeMessageRequest{}
	if err := s.getRequest(conn, challengeRequest); err != nil {
		return fmt.Errorf("get request: %w\n", err)
	}
	return nil
}

func (s *Server) writeChallengeResponse(conn net.Conn, challenge []byte, difficulty int) error {
	challengeResponse := &messages.ChallengeMessageResponse{
		Type:       messages.ChallengeResponse,
		Challenge:  challenge,
		Difficulty: difficulty,
	}
	if err := s.sendResponse(conn, challengeResponse); err != nil {
		return fmt.Errorf("send response: %w\n", err)
	}
	return nil
}

func (s *Server) readSolutionRequest(conn net.Conn) ([]byte, error) {
	solutionRequest := &messages.SolutionMessageRequest{
		Solution: make([]byte, 0),
	}
	if err := s.getRequest(conn, solutionRequest); err != nil {
		return nil, fmt.Errorf("get request: %w\n", err)
	}
	return solutionRequest.Solution, nil
}

func (s *Server) handleError(conn net.Conn, err error, content string) {
	errorMessage := &messages.ErrorMessage{
		Type:         messages.Error,
		ErrorMessage: content,
	}
	if errResp := s.sendResponse(conn, errorMessage); errResp != nil {
		slog.With("error", errResp.Error()).Error("send error response")
	}
	slog.With("error", err.Error()).Error(content)
}

func (s *Server) getRequest(conn net.Conn, req interface{}) error {
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(req); err != nil {
		return fmt.Errorf("decode request: %w\n", err)
	}

	if err := validator.GetInstance().Struct(req); err != nil {
		return fmt.Errorf("validate data: %w\n", err)
	}

	slog.With("remote_address", conn.RemoteAddr().String()).Info("tcp connection")

	return nil
}

func (s *Server) sendResponse(conn net.Conn, data interface{}) error {
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode data: %w\n", err)
	}
	slog.With("remote_address", conn.RemoteAddr().String()).Info("tcp connection wrote")
	return nil
}
