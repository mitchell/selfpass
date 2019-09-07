package crypto

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	PBKDF2Rounds = 4096
	KeyLength    = 32
)

func GeneratePBKDF2Key(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, PBKDF2Rounds, KeyLength, sha256.New)
}
