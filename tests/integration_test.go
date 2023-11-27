package tests

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/to77e/word-of-wisdom/internal/client"
	"github.com/to77e/word-of-wisdom/internal/proofofwork"
)

const (
	defaultAddress = "localhost:11001"
)

type SuiteServer struct {
	suite.Suite
	ctx       context.Context
	validator *validator.Validate
}

func TestIntegrationServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SuiteServer))
}

func (s *SuiteServer) SetupSuite() {
	s.ctx = context.Background()
	s.validator = validator.New(validator.WithRequiredStructEnabled())
}

func (s *SuiteServer) TearDownSuite() {
	s.ctx.Done()
}

func (s *SuiteServer) TestServer() {
	t := s.T()
	t.Parallel()

	s.Run("get word of wisdom quote", func() {
		tcpClient, err := client.New(defaultAddress)
		assert.NoError(t, err)
		defer tcpClient.Close() //nolint: errcheck

		err = tcpClient.WriteChallengeRequest()
		assert.NoError(t, err)

		challenge, err := tcpClient.ReadChallengeResponse()
		assert.NoError(t, err)

		pow := proofofwork.New()
		pow.SetChallenge(challenge)

		err = pow.ComputeSolution()
		assert.NoError(t, err)

		err = tcpClient.WriteSolutionRequest(pow.GetSolution())
		assert.NoError(t, err)

		result, err := tcpClient.ReadSolutionResponse()
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	s.Run("compute wrong solution", func() {
		tcpClient, err := client.New(defaultAddress)
		assert.NoError(t, err)
		defer tcpClient.Close() //nolint: errcheck

		err = tcpClient.WriteChallengeRequest()
		assert.NoError(t, err)

		challenge, err := tcpClient.ReadChallengeResponse()
		assert.NoError(t, err)

		pow := proofofwork.New()
		pow.SetChallenge(challenge)
		pow.SetSolution([]byte("wrong solution"))
		err = tcpClient.WriteSolutionRequest(pow.GetSolution())
		assert.NoError(t, err)

		result, err := tcpClient.ReadSolutionResponse()
		assert.Error(t, err, "solution is not valid")
		assert.Empty(t, result)
	})

	s.Run("not send start request", func() {
		tcpClient, err := client.New(defaultAddress)
		assert.NoError(t, err)
		defer tcpClient.Close() //nolint: errcheck

		pow := proofofwork.New()
		pow.SetChallenge([]byte("wrong challenge"))

		err = pow.ComputeSolution()
		assert.NoError(t, err)

		err = tcpClient.WriteSolutionRequest(pow.GetSolution())
		assert.NoError(t, err)

		result, err := tcpClient.ReadSolutionResponse()
		assert.Error(t, err, "solution is not valid")
		assert.Empty(t, result)
	})

	s.Run("not send solution request", func() {
		tcpClient, err := client.New(defaultAddress)
		assert.NoError(t, err)
		defer tcpClient.Close() //nolint: errcheck

		err = tcpClient.WriteChallengeRequest()
		assert.NoError(t, err)

		challenge, err := tcpClient.ReadChallengeResponse()
		assert.NoError(t, err)

		pow := proofofwork.New()
		pow.SetChallenge(challenge)

		err = pow.ComputeSolution()
		assert.NoError(t, err)

		result, err := tcpClient.ReadSolutionResponse()
		assert.Error(t, err, "solution is not valid")
		assert.Empty(t, result)
	})
}
