package companydto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

type CompanyUpdateDTO struct {
	TradeName   *string                      `json:"trade_name"`
	Cnpj        *string                      `json:"cnpj"`
	Email       *string                      `json:"email"`
	Contacts    []string                     `json:"contacts"`
	Address     *addressdto.AddressUpdateDTO `json:"address"`
	Preferences companyentity.Preferences    `json:"preferences"`
	// Fiscal fields
	FiscalEnabled     *bool   `json:"fiscal_enabled,omitempty"`
	InscricaoEstadual *string `json:"inscricao_estadual,omitempty"`
	RegimeTributario  *int    `json:"regime_tributario,omitempty"`
	CNAE              *string `json:"cnae,omitempty"`
	CRT               *int    `json:"crt,omitempty"`
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

	if c.Address != nil {
		if err := c.Address.UpdateDomain(company.Address); err != nil {
			return err
		}
	}

	// Update fiscal fields if provided
	if c.FiscalEnabled != nil {
		company.FiscalEnabled = *c.FiscalEnabled
	}
	if c.InscricaoEstadual != nil {
		company.InscricaoEstadual = *c.InscricaoEstadual
	}
	if c.RegimeTributario != nil {
		company.RegimeTributario = *c.RegimeTributario
	}
	if c.CNAE != nil {
		company.CNAE = *c.CNAE
	}
	if c.CRT != nil {
		company.CRT = *c.CRT
	}

	return nil
}
