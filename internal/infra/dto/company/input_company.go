package companydto

import (
	"errors"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

var (
	ErrMustBeCNPJ     = errors.New("cnpj is required")
	ErrMustBeContacts = errors.New("contacts is required")
)

type CompanyInput struct {
	companyentity.CompanyCommonAttributes
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
