package companydto

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
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
}

func (c *CompanyDTO) FromDomain(company *companyentity.Company) {
	if company == nil {
		return
	}
	*c = CompanyDTO{
		ID:           company.ID,
		SchemaName:   company.SchemaName,
		BusinessName: company.BusinessName,
		TradeName:    company.TradeName,
		Cnpj:         company.Cnpj,
		Email:        company.Email,
		Contacts:     company.Contacts,
		Address:      &addressdto.AddressDTO{},
		Users:        []UserDTO{},
		Preferences:  company.Preferences,
		IsBlocked:    company.IsBlocked,
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
}
