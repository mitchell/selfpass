package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/cloudflare/redoctober/padding"
)

func GenerateKeyFromPassword(pass []byte) ([]byte, error) {
	if len(pass) < 8 {
		return nil, fmt.Errorf("master password must be at least 8 characters")
	}

	for idx := 0; len(pass) < 32; idx++ {
		pass = append(pass, pass[idx])

		if idx == len(pass) {
			idx = 0
		}
	}

	return pass, nil
}

func CombinePasswordAndKey(pass, key []byte) ([]byte, error) {
	if len(pass) < 8 {
		return nil, fmt.Errorf("master password must be at least 8 characters")
	}
	if len(key) != 16 {
		return nil, fmt.Errorf("key was not of length 16")
	}

	for idx := 0; len(pass) < 16; idx++ {
		pass = append(pass, pass[idx])
	}

	return append(pass[:16], key...), nil
}

func CBCEncrypt(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key is not 32 bytes")
	}

	plaintext = padding.AddPadding(plaintext)

	if len(plaintext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func CBCDecrypt(key []byte, ciphertext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key is not 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]

	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return padding.RemovePadding(ciphertext)
}
