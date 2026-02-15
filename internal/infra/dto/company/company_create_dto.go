package companydto

import (
	"errors"
)

var (
	ErrMustBeCNPJ       = errors.New("cnpj is required")
	ErrMustBeContacts   = errors.New("contacts is required")
	ErrMustBeCategoryID = errors.New("category_id is required")
)

type CompanyCreateDTO struct {
	TradeName   string   `json:"trade_name"`
	Cnpj        string   `json:"cnpj"`
	Contacts    []string `json:"contacts"`
	CategoryIDs []string `json:"category_ids"`
}

func (c *CompanyCreateDTO) validate() error {
	if c.Cnpj == "" {
		return ErrMustBeCNPJ
	}

	if len(c.Contacts) == 0 {
		return ErrMustBeContacts
	}

	if len(c.CategoryIDs) == 0 {
		return ErrMustBeCategoryID
	}

	return nil
}

func (c *CompanyCreateDTO) ToDomain() (cnpj string, tradeName string, contacts []string, categoryIDs []string, err error) {
	if err := c.validate(); err != nil {
		return "", "", nil, nil, err
	}

	return c.Cnpj, c.TradeName, c.Contacts, c.CategoryIDs, nil
}
