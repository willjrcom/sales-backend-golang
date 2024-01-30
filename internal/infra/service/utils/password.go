package utils

import "regexp"

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	if !containsUppercase(password) {
		return false
	}

	if !containsLowercase(password) {
		return false
	}

	if !containsDigit(password) {
		return false
	}

	if !containsSymbol(password) {
		return false
	}

	return true
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
