package companydto

import (
	"time"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	companycategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company_category"
)

type CompanyDTO struct {
	ID           uuid.UUID                 `json:"id"`
	SchemaName   string                    `json:"schema_name"`
	BusinessName string                    `json:"business_name"`
	TradeName    string                    `json:"trade_name"`
	Cnpj         string                    `json:"cnpj"`
	Email        string                    `json:"email"`
	Contacts     []string                  `json:"contacts"`
	Address      *addressdto.AddressDTO    `json:"address,omitempty"`
	Users        []UserDTO                 `json:"users,omitempty"`
	Preferences  companyentity.Preferences `json:"preferences,omitempty"`
	IsBlocked    bool                      `json:"is_blocked,omitempty"`

	// Categories
	Categories []companycategorydto.CompanyCategoryDTO `json:"categories,omitempty"`

	MonthlyPaymentDueDay          int        `json:"monthly_payment_due_day,omitempty"`
	MonthlyPaymentDueDayUpdatedAt *time.Time `json:"monthly_payment_due_day_updated_at,omitempty"`
}

func (c *CompanyDTO) FromDomain(company *companyentity.Company) {
	if company == nil {
		return
	}
	*c = CompanyDTO{
		ID:                            company.ID,
		SchemaName:                    company.SchemaName,
		BusinessName:                  company.BusinessName,
		TradeName:                     company.TradeName,
		Cnpj:                          company.Cnpj,
		Email:                         company.Email,
		Contacts:                      company.Contacts,
		Address:                       &addressdto.AddressDTO{},
		Users:                         []UserDTO{},
		Preferences:                   company.Preferences,
		IsBlocked:                     company.IsBlocked,
		Categories:                    []companycategorydto.CompanyCategoryDTO{},
		MonthlyPaymentDueDay:          company.MonthlyPaymentDueDay,
		MonthlyPaymentDueDayUpdatedAt: company.MonthlyPaymentDueDayUpdatedAt,
	}

	c.Address.FromDomain(company.Address)

	for _, user := range company.Users {
		userDTO := UserDTO{}
		userDTO.FromDomain(&user)
		c.Users = append(c.Users, userDTO)
	}

	if company.Address == nil {
		c.Address = nil
	}

	for _, category := range company.Categories {
		categoryDTO := companycategorydto.CompanyCategoryDTO{}
		categoryDTO.FromDomain(&category)
		c.Categories = append(c.Categories, categoryDTO)
	}
}
