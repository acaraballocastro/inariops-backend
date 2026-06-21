package auth

import (
	"crypto/rand"
	"math/big"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CheckNewPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true

		case unicode.IsLower(char):
			hasLower = true

		case unicode.IsDigit(char):
			hasDigit = true

		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func createDefaultPassword() (string, error) {
	const passwordLength = 12

	const (
		upper   = "ABCDEFGHJKLMNPQRSTUVWXYZ"
		lower   = "abcdefghijkmnopqrstuvwxyz"
		digits  = "23456789"
		symbols = "!@#$%^&*()-_=+"
	)

	allChars := upper + lower + digits + symbols

	password := make([]byte, passwordLength)

	required := []string{
		upper,
		lower,
		digits,
		symbols,
	}

	// Garantizar requisitos mínimos
	for i, charset := range required {
		c, err := randomChar(charset)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Completar el resto
	for i := len(required); i < passwordLength; i++ {
		c, err := randomChar(allChars)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Fisher-Yates
	for i := len(password) - 1; i > 0; i-- {
		j, err := cryptoRandInt(i + 1)
		if err != nil {
			return "", err
		}

		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

func randomChar(charset string) (byte, error) {
	idx, err := cryptoRandInt(len(charset))
	if err != nil {
		return 0, err
	}

	return charset[idx], nil
}

func cryptoRandInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}
