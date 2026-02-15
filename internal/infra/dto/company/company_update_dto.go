package companydto

import (
	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	companycategoryentity "github.com/willjrcom/sales-backend-go/internal/domain/company_category"
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
	CategoryIDs  []string                     `json:"category_ids"`

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

	if len(c.CategoryIDs) > 0 {
		company.Categories = []companycategoryentity.CompanyCategory{}
		for _, id := range c.CategoryIDs {
			categoryIDUUID, err := uuid.Parse(id)
			if err != nil {
				return err
			}

			category := companycategoryentity.CompanyCategory{}
			category.ID = categoryIDUUID
			company.Categories = append(company.Categories, category)
		}
	}

	return nil
}
