package proofofwork

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/to77e/word-of-wisdom/internal/scrypt"
)

const (
	defaultDifficulty = 1
	defaultBufferSize = 32
)

type ProofOfWork struct {
	challenge []byte
	solution  []byte
}

func New() *ProofOfWork {
	return &ProofOfWork{
		challenge: make([]byte, defaultBufferSize),
		solution:  make([]byte, defaultBufferSize),
	}
}

func (p *ProofOfWork) GenerateChallenge() error {
	if _, err := io.ReadFull(rand.Reader, p.challenge); err != nil {
		return fmt.Errorf("read challenge: %w\n", err)
	}
	return nil
}

func (p *ProofOfWork) GetChallenge() []byte {
	return p.challenge
}

func (p *ProofOfWork) SetChallenge(challenge []byte) {
	p.challenge = challenge
}

func (p *ProofOfWork) GetSolution() []byte {
	return p.solution
}

func (p *ProofOfWork) SetSolution(solution []byte) {
	p.solution = solution
}

func (p *ProofOfWork) VerifySolution() (bool, error) {
	hash, err := scrypt.GenerateHash(p.solution, p.challenge)
	if err != nil {
		return false, fmt.Errorf("scrypt key: %w\n", err)
	}
	return bytes.HasPrefix(hash, bytes.Repeat([]byte{0}, defaultDifficulty)), nil
}

func (p *ProofOfWork) ComputeSolution() error {
	for {
		if _, err := rand.Read(p.solution); err != nil {
			return fmt.Errorf("read solution: %w\n", err)
		}

		hash, err := scrypt.GenerateHash(p.solution, p.challenge)
		if err != nil {
			return fmt.Errorf("generate hash: %w\n", err)
		}

		if bytes.HasPrefix(hash, bytes.Repeat([]byte{0}, defaultDifficulty)) {
			return nil
		}
	}
}
