package companydto

import (
	"time"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

type CompanyDTO struct {
	ID                    uuid.UUID                 `json:"id"`
	SchemaName            string                    `json:"schema_name"`
	BusinessName          string                    `json:"business_name"`
	TradeName             string                    `json:"trade_name"`
	Cnpj                  string                    `json:"cnpj"`
	Email                 string                    `json:"email"`
	Contacts              []string                  `json:"contacts"`
	Address               *addressdto.AddressDTO    `json:"address,omitempty"`
	Users                 []UserDTO                 `json:"users,omitempty"`
	Preferences           companyentity.Preferences `json:"preferences,omitempty"`
	IsBlocked             bool                      `json:"is_blocked,omitempty"`
	SubscriptionExpiresAt *time.Time                `json:"subscription_expires_at,omitempty"`
	// Fiscal fields
	FiscalEnabled     bool   `json:"fiscal_enabled,omitempty"`
	InscricaoEstadual string `json:"inscricao_estadual,omitempty"`
	RegimeTributario  int    `json:"regime_tributario,omitempty"`
	CNAE              string `json:"cnae,omitempty"`
	CRT               int    `json:"crt,omitempty"`
}

func (c *CompanyDTO) FromDomain(company *companyentity.Company) {
	if company == nil {
		return
	}
	*c = CompanyDTO{
		ID:                    company.ID,
		SchemaName:            company.SchemaName,
		BusinessName:          company.BusinessName,
		TradeName:             company.TradeName,
		Cnpj:                  company.Cnpj,
		Email:                 company.Email,
		Contacts:              company.Contacts,
		Address:               &addressdto.AddressDTO{},
		Users:                 []UserDTO{},
		Preferences:           company.Preferences,
		IsBlocked:             company.IsBlocked,
		SubscriptionExpiresAt: company.SubscriptionExpiresAt,
		FiscalEnabled:         company.FiscalEnabled,
		InscricaoEstadual:     company.InscricaoEstadual,
		RegimeTributario:      company.RegimeTributario,
		CNAE:                  company.CNAE,
		CRT:                   company.CRT,
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
