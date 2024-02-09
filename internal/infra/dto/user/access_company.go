package userdto

import (
	"errors"
)

var (
	ErrSchemaRequired = errors.New("schema required")
)

type AccessCompanyInput struct {
	Schema *string `json:"schema"`
}

func (u *AccessCompanyInput) validate() error {
	if u.Schema == nil {
		return ErrSchemaRequired
	}
	return nil
}

func (u *AccessCompanyInput) ToModel() (*string, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	return u.Schema, nil
}
