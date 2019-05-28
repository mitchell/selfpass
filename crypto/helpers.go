package crypto

import (
	"fmt"
	"math/rand"
	"time"
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

func GeneratePassword(length int, numbers, specials bool) string {
	const alphas = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const alphanumerics = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const alphasAndSpecials = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()"
	const alphanumericsAndSpecials = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()"
	rand.Seed(time.Now().UnixNano())
	pass := make([]byte, length)

	switch {
	case numbers && specials:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphanumericsAndSpecials[rand.Int63()%int64(len(alphanumericsAndSpecials))]
		}
	case numbers:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphanumerics[rand.Int63()%int64(len(alphanumerics))]
		}
	case specials:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphasAndSpecials[rand.Int63()%int64(len(alphasAndSpecials))]
		}
	default:
		for idx := 0; idx < length; idx++ {
			pass[idx] = alphas[rand.Int63()%int64(len(alphas))]
		}
	}

	return string(pass)
}
