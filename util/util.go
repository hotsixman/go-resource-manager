package util

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
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

func GetParentDirectory(path string) (string, error) {
	if path[0] == '/' {
		return "", errors.New("경로는 '/'으로 시작해야합니다")
	}
	if path == "/" {
		return "", errors.New("경로는 \"/\"일 수 없습니다")
	}

	parentPath := ""
	names := strings.Split(path, "/")
	namesLen := len(names)
	for i, name := range names {
		if i == 0 {
			continue
		}
		if i == namesLen-1 {
			continue
		}
		parentPath += "/" + name
	}

	return parentPath, nil
}
