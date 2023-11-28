package proofofwork

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

const (
	bufferSize = 32
)

type ProofOfWork struct {
	challenge  []byte
	solution   []byte
	difficulty int
}

func New(difficulty int) *ProofOfWork {
	return &ProofOfWork{
		challenge:  make([]byte, bufferSize),
		solution:   make([]byte, bufferSize),
		difficulty: difficulty,
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

func (p *ProofOfWork) GetDifficulty() int {
	return p.difficulty
}

func (p *ProofOfWork) VerifySolution() bool {
	hash := sha256.Sum256(p.solution)
	return bytes.HasPrefix(hash[:], bytes.Repeat([]byte{0}, p.difficulty))
}

func (p *ProofOfWork) ComputeSolution() error {
	for {
		if _, err := rand.Read(p.solution); err != nil {
			return fmt.Errorf("read solution: %w\n", err)
		}
		hash := sha256.Sum256(p.solution)
		if bytes.HasPrefix(hash[:], bytes.Repeat([]byte{0}, p.difficulty)) {
			return nil
		}
	}
}
