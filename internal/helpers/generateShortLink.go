package helpers

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortLink(length int) string {
	if length < 0 {
		return ""
	}

	result := make([]byte, length)

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
