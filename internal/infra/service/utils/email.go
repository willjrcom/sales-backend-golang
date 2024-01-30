package utils

import "regexp"

func IsEmailValid(email string) bool {
	// Expressão regular para validar endereços de e-mail simples
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
