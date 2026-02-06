package companydto

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

type CompanyUpdateDTO struct {
	BusinessName *string                      `json:"business_name"`
	TradeName    *string                      `json:"trade_name"`
	Cnpj         *string                      `json:"cnpj"`
	Email        *string                      `json:"email"`
	Contacts     []string                     `json:"contacts"`
	Address      *addressdto.AddressUpdateDTO `json:"address"`
	Preferences  companyentity.Preferences    `json:"preferences"`

	MonthlyPaymentDueDay *int `json:"monthly_payment_due_day,omitempty"`
}

func (c *CompanyUpdateDTO) validate() error {
	return nil
}

func (c *CompanyUpdateDTO) UpdateDomain(company *companyentity.Company) (err error) {
	if err := c.validate(); err != nil {
		return err
	}

	if c.BusinessName != nil {
		company.BusinessName = *c.BusinessName
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

	if c.Address != nil {
		if company.Address == nil {
			company.AddAddress(&addressentity.AddressCommonAttributes{})
		}

		if err := c.Address.UpdateDomain(company.Address); err != nil {
			return err
		}
	}

	if c.MonthlyPaymentDueDay != nil {
		company.MonthlyPaymentDueDay = *c.MonthlyPaymentDueDay
	}

	return nil
}
