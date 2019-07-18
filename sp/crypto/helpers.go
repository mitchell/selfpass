package crypto

import (
	"fmt"
	"math/rand"
	"time"
)

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
	const alphaValues = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numberValues = "1234567890"
	const specialValues = "!@#$%^&*()_-+="
	rand.Seed(time.Now().UnixNano())
	pass := make([]byte, length)

	values := alphaValues

	switch {
	case numbers && specials:
		values += numberValues + specialValues
	case numbers:
		values += numberValues
	case specials:
		values += specialValues
	}

	for idx := range pass {
		pass[idx] = values[rand.Int63()%int64(len(values))]
	}

	return string(pass)
}
