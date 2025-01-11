package companydto

import (
	"errors"
)

var (
	ErrCpfRequired = errors.New("cpf is required")
)

type UserSearchDTO struct {
	Cpf string `json:"cpf"`
}

func (u *UserSearchDTO) validate() error {
	if u.Cpf == "" {
		return ErrCpfRequired
	}

	return nil
}

func (u *UserSearchDTO) ToDomain() (string, error) {
	if err := u.validate(); err != nil {
		return "", err
	}

	return u.Cpf, nil
}
