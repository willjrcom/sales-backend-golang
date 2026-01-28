package fiscalsettingsdto

import (
	fiscalsettingsentity "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_settings"
)

type FiscalSettingsDTO struct {
	CompanyRegistryID     int64  `json:"company_registry_id,omitempty"`
	FiscalEnabled         bool   `json:"fiscal_enabled"`
	StateRegistration     string `json:"state_registration"`
	TaxRegime             int    `json:"tax_regime"`
	CNAE                  string `json:"cnae"`
	CRT                   int    `json:"crt"`
	IsSimpleNational      bool   `json:"is_simple_national"`
	MunicipalRegistration string `json:"municipal_registration"`

	// Preferences
	ShowTaxBreakdown     bool `json:"show_tax_breakdown"`
	SendEmailToRecipient bool `json:"send_email_to_recipient"`

	// Company Identity
	BusinessName string `json:"business_name"` // Razão Social
	TradeName    string `json:"trade_name"`    // Nome Fantasia
	Cnpj         string `json:"cnpj"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`

	// Address
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	UF           string `json:"uf"`
	Cep          string `json:"cep"`
}

func (d *FiscalSettingsDTO) FromDomain(entity *fiscalsettingsentity.FiscalSettings) {
	d.CompanyRegistryID = entity.CompanyRegistryID
	d.FiscalEnabled = entity.IsActive
	d.StateRegistration = entity.StateRegistration
	d.TaxRegime = entity.TaxRegime
	d.CNAE = entity.CNAE
	d.CRT = entity.CRT
	d.IsSimpleNational = entity.IsSimpleNational
	d.MunicipalRegistration = entity.MunicipalRegistration
	d.ShowTaxBreakdown = entity.ShowTaxBreakdown
	d.SendEmailToRecipient = entity.SendEmailToRecipient

	d.BusinessName = entity.BusinessName
	d.TradeName = entity.TradeName
	d.Cnpj = entity.Cnpj
	d.Email = entity.Email
	d.Phone = entity.Phone

	d.Street = entity.Address.Street
	d.Number = entity.Address.Number
	d.Complement = entity.Address.Complement
	d.Neighborhood = entity.Address.Neighborhood
	d.City = entity.Address.City
	d.UF = entity.Address.UF
	d.Cep = entity.Address.Cep
}
