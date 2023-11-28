package proofofwork

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_VerifySolution(t *testing.T) {
	const (
		difficulty = 1
	)
	type fields struct {
		challenge func(*ProofOfWork) []byte
		solution  func(*ProofOfWork) []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive",
			fields: fields{
				challenge: func(pow *ProofOfWork) []byte {
					_ = pow.GenerateChallenge() //nolint:errcheck
					return pow.GetChallenge()
				},
				solution: func(pow *ProofOfWork) []byte {
					_ = pow.ComputeSolution() //nolint:errcheck
					return pow.GetSolution()
				},
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "negative",
			fields: fields{
				challenge: func(pow *ProofOfWork) []byte {
					_ = pow.GenerateChallenge() //nolint:errcheck
					return pow.GetChallenge()
				},
				solution: func(pow *ProofOfWork) []byte {
					return []byte("wrong solution")
				},
			},
			want: false,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return fmt.Errorf("scrypt key: %w\n", err) != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(difficulty)
			p.challenge = tt.fields.challenge(p)
			p.solution = tt.fields.solution(p)

			got := p.VerifySolution()
			assert.Equalf(t, tt.want, got, "VerifySolution() = %v, want %v", got, tt.want)
		})
	}
}
