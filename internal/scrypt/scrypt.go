package scrypt

import (
	"fmt"

	"golang.org/x/crypto/scrypt"
)

var (
	N       = 4  // CPU/memory cost factor
	r       = 2  // block size
	p       = 1  // parallelization factor
	keySize = 32 // desired key length
)

func GenerateHash(solution, challenge []byte) ([]byte, error) {
	hash, err := scrypt.Key(solution, challenge, N, r, p, keySize)
	if err != nil {
		return nil, fmt.Errorf("scrypt key: %w\n", err)
	}
	return hash, nil
}
