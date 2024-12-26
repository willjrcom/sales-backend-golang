package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !containsUppercase(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !containsLowercase(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !containsDigit(password) {
		return errors.New("password must contain at least one digit")
	}

	if !containsSymbol(password) {
		return errors.New("password must contain at least one symbol")
	}

	return nil
}

func containsUppercase(s string) bool {
	uppercaseRegex := `[A-Z]`
	match, _ := regexp.MatchString(uppercaseRegex, s)
	return match
}

func containsLowercase(s string) bool {
	lowercaseRegex := `[a-z]`
	match, _ := regexp.MatchString(lowercaseRegex, s)
	return match
}

func containsDigit(s string) bool {
	digitRegex := `[0-9]`
	match, _ := regexp.MatchString(digitRegex, s)
	return match
}

func containsSymbol(s string) bool {
	symbolRegex := `[!@#$%^&*()-=_+{}\[\]:;<>,.?/~]`
	match, _ := regexp.MatchString(symbolRegex, s)
	return match
}

func GeneratePassword(length int, useUppercase bool, useNumbers bool, useSymbols bool) string {
	lowercaseChars := "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars := "0123456789"
	symbolChars := "!@#$%^&*()_+-=[]{}|;:'\",.<>?/`~"

	// Montar conjunto de caracteres
	characters := lowercaseChars
	if useUppercase {
		characters += uppercaseChars
	}
	if useNumbers {
		characters += numberChars
	}
	if useSymbols {
		characters += symbolChars
	}

	if len(characters) == 0 {
		return ""
	}

	// Gerar a senha
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		password[i] = characters[randomIndex.Int64()]
	}

	return string(password)
}
