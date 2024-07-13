package crypto

import (
	"github.com/gagraler/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var log = logger.SugaredLogger()

func Encryption(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err.Error())
	}

	return string(hash)
}

func Compare(ciphertext, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(ciphertext), []byte(password))
	if err != nil {
		log.Error(err.Error())
		return false
	}

	return true
}
