package auth

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func createDefaultPassword() (string, error) {
	const passwordLength = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/~`"

	password := make([]byte, passwordLength)
	for i := range password {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			return "", err
		}
		password[i] = charset[int(randomByte[0])%len(charset)]
	}

	return string(password), nil
}
