package utils

import (
	"golang.org/x/crypto/bcrypt"

	"fmt"
)

func HashPassword(Password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to generate password: %w", err)
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedpassword string, password string) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedpassword),
		[]byte(password),
	)

	if err != nil {
		return false
	}

	return true
}
