package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashB, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate password: %w", err)
	}
	return string(hashB), nil
}

func CheckPassword(password, hashB string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashB), []byte(password))
}
