package util

import (
	"crypto/rand"
	"math/big"
)

func CloneSlice[T any](original []T) []T {
	clone := make([]T, len(original))
	copy(clone, original)
	return clone
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		// crypto/rand.Intn은 보안적으로 안전한 난수를 반환합니다.
		var num *big.Int
		var err error
		for {
			num, err = rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err == nil {
				break
			}
		}
		result[i] = charset[num.Int64()]
	}

	return string(result)
}
