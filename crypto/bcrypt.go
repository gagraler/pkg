package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Encryption(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword err: %v", err)
	}

	return string(hash), err
}

func Compare(ciphertext, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(ciphertext), []byte(password))
	if err != nil {
		return false, fmt.Errorf("bcrypt.CompareHashAndPassword err: %v", err)
	}

	return true, nil
}
