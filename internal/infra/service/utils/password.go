package utils

import (
	"errors"
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
