package companydto

import (
	"errors"
)

var (
	ErrMustBeEmail = errors.New("email is required")
)

type UserToCompanyDTO struct {
	Email string `json:"email"`
}

func (c *UserToCompanyDTO) validate() error {
	if c.Email == "" {
		return ErrMustBeEmail
	}
	return nil
}

func (u *UserToCompanyDTO) ToDomain() (string, error) {
	if err := u.validate(); err != nil {
		return "", err
	}

	return u.Email, nil
}
