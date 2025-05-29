package companydto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyUpdateDTO struct {
	TradeName   *string                   `json:"trade_name"`
	Cnpj        *string                   `json:"cnpj"`
	Email       *string                   `json:"email"`
	Contacts    []string                  `json:"contacts"`
	Preferences companyentity.Preferences `json:"preferences"`
}

func (c *CompanyUpdateDTO) validate() error {
	return nil
}

func (c *CompanyUpdateDTO) UpdateDomain(company *companyentity.Company) (err error) {
	if err := c.validate(); err != nil {
		return err
	}

	if c.TradeName != nil {
		company.TradeName = *c.TradeName
	}
	if c.Cnpj != nil {
		company.Cnpj = *c.Cnpj
	}
	if c.Email != nil {
		company.Email = *c.Email
	}
	if len(c.Contacts) > 0 {
		company.Contacts = c.Contacts
	}
	if c.Preferences != nil {
		company.Preferences = c.Preferences
	}

	return nil
}
