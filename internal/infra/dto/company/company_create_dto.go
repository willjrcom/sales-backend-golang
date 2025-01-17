package companydto

import (
	"errors"
)

var (
	ErrMustBeCNPJ     = errors.New("cnpj is required")
	ErrMustBeContacts = errors.New("contacts is required")
)

type CompanyCreateDTO struct {
	TradeName string   `json:"trade_name"`
	Cnpj      string   `json:"cnpj"`
	Contacts  []string `json:"contacts"`
}

func (c *CompanyCreateDTO) validate() error {
	if c.Cnpj == "" {
		return ErrMustBeCNPJ
	}

	if len(c.Contacts) == 0 {
		return ErrMustBeContacts
	}
	return nil
}

func (c *CompanyCreateDTO) ToDomain() (cnpj string, tradeName string, contacts []string, err error) {
	if err := c.validate(); err != nil {
		return "", "", nil, err
	}

	return c.Cnpj, c.TradeName, c.Contacts, nil
}
