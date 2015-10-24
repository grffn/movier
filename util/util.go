package util

import "golang.org/x/crypto/scrypt"

//GeneratePassword genarates encrypted password
func GeneratePassword(password []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 16384, 8, 1, 32)
}
