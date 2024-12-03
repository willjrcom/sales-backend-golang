package companydto

import (
	"errors"
)

var (
	ErrMustBeCNPJ     = errors.New("cnpj is required")
	ErrMustBeContacts = errors.New("contacts is required")
)

type CompanyInput struct {
	TradeName string   `json:"trade_name"`
	Cnpj      string   `json:"cnpj"`
	Email     string   `json:"email"`
	Contacts  []string `json:"contacts"`
}

func (c *CompanyInput) validate() error {
	if c.Cnpj == "" {
		return ErrMustBeCNPJ
	}

	if len(c.Contacts) == 0 {
		return ErrMustBeContacts
	}
	return nil
}

func (c *CompanyInput) ToModel() (cnpj string, tradeName string, email string, contacts []string, err error) {
	if err := c.validate(); err != nil {
		return "", "", "", nil, err
	}

	return c.Cnpj, c.TradeName, c.Email, c.Contacts, nil
}
