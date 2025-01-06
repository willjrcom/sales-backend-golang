package userdto

import (
	"errors"
)

var (
	ErrSchemaRequired = errors.New("schema required")
)

type UserSchemaDTO struct {
	Schema *string `json:"schema"`
}

func (u *UserSchemaDTO) validate() error {
	if u.Schema == nil {
		return ErrSchemaRequired
	}
	return nil
}

func (u *UserSchemaDTO) ToDomain() (*string, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	return u.Schema, nil
}
