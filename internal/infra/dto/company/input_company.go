package companydto

import (
	"errors"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrMustBeName        = errors.New("name is required")
	ErrMustBeCNPJ        = errors.New("cnpj is required")
	ErrMustBeContacts    = errors.New("contacts is required")
	ErrSchemaNameIsEmpty = errors.New("schema name is required")
)

type CompanyInput struct {
	companyentity.CompanyCommonAttributes
}

func (c *CompanyInput) validate() error {
	if c.Name == "" {
		return ErrMustBeName
	}

	if c.Cnpj == "" {
		return ErrMustBeCNPJ
	}

	if err := c.Address.Validate(); err != nil {
		return err
	}

	if len(c.Contacts) == 0 {
		return ErrMustBeContacts
	}

	if c.SchemaName == "" {
		return ErrSchemaNameIsEmpty
	}
	return nil
}

func (c *CompanyInput) ToModel() (*companyentity.Company, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}
	return &companyentity.Company{
		Entity:                  entity.NewEntity(),
		CompanyCommonAttributes: c.CompanyCommonAttributes,
	}, nil
}
