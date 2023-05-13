package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)


func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hash), err
}

func Verify(hashed string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func PairHash(message []byte, secret string) string {

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	hash := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)

}