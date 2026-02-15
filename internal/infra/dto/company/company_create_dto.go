package companydto

import (
	"errors"

	"github.com/google/uuid"
	companycategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company_category"
)

var (
	ErrMustBeCNPJ       = errors.New("cnpj is required")
	ErrMustBeContacts   = errors.New("contacts is required")
	ErrMustBeCategoryID = errors.New("category_id is required")
)

type CompanyCreateDTO struct {
	TradeName  string                                  `json:"trade_name"`
	Cnpj       string                                  `json:"cnpj"`
	Contacts   []string                                `json:"contacts"`
	Categories []companycategorydto.CompanyCategoryDTO `json:"categories"`
}

func (c *CompanyCreateDTO) validate() error {
	if c.Cnpj == "" {
		return ErrMustBeCNPJ
	}

	if len(c.Contacts) == 0 {
		return ErrMustBeContacts
	}

	if len(c.Categories) == 0 {
		return ErrMustBeCategoryID
	}

	return nil
}

func (c *CompanyCreateDTO) ToDomain() (cnpj string, tradeName string, contacts []string, categoryIDs []string, err error) {
	if err := c.validate(); err != nil {
		return "", "", nil, nil, err
	}

	for _, category := range c.Categories {
		if category.ID != uuid.Nil {
			categoryIDs = append(categoryIDs, category.ID.String())
		}
	}

	return c.Cnpj, c.TradeName, c.Contacts, categoryIDs, nil
}
