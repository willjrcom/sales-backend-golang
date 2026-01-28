package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	fiscalsettingsentity "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_settings"
)

type FiscalSettings struct {
	bun.BaseModel `bun:"table:fiscal_settings,alias:fs"`

	ID                    uuid.UUID `bun:"id,pk,type:uuid"`
	CompanyID             uuid.UUID `bun:"company_id,type:uuid,notnull"`
	CompanyRegistryID     int64     `bun:"company_registry_id"`
	IsActive              bool      `bun:"is_active"`
	StateRegistration     string    `bun:"state_registration"`     // InscricaoEstadual
	TaxRegime             int       `bun:"tax_regime"`             // RegimeTributario
	CNAE                  string    `bun:"cnae"`                   // CNAE
	CRT                   int       `bun:"crt"`                    // CRT
	IsSimpleNational      bool      `bun:"is_simple_national"`     // SimplesNacional
	MunicipalRegistration string    `bun:"municipal_registration"` // InscricaoMunicipal

	// Preferences
	ShowTaxBreakdown     bool `bun:"show_tax_breakdown"`      // DiscriminaImpostos
	SendEmailToRecipient bool `bun:"send_email_to_recipient"` // EnviarEmailDestinatario

	// Company Identity
	BusinessName string `bun:"business_name"`
	TradeName    string `bun:"trade_name"`
	Cnpj         string `bun:"cnpj"`
	Email        string `bun:"email"`
	Phone        string `bun:"phone"`

	// Address Fields
	Street       string `bun:"street"`
	Number       string `bun:"number"`
	Complement   string `bun:"complement"`
	Neighborhood string `bun:"neighborhood"`
	City         string `bun:"city"`
	UF           string `bun:"uf"`
	Cep          string `bun:"cep"`

	CreatedAt time.Time  `bun:"created_at,default:current_timestamp"`
	UpdatedAt *time.Time `bun:"updated_at"`
}

func (m *FiscalSettings) ToDomain() *fiscalsettingsentity.FiscalSettings {
	var updatedAt time.Time
	if m.UpdatedAt != nil {
		updatedAt = *m.UpdatedAt
	}

	return &fiscalsettingsentity.FiscalSettings{
		Entity: entity.Entity{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: updatedAt,
		},
		CompanyID:             m.CompanyID,
		CompanyRegistryID:     m.CompanyRegistryID,
		IsActive:              m.IsActive,
		StateRegistration:     m.StateRegistration,
		TaxRegime:             m.TaxRegime,
		CNAE:                  m.CNAE,
		CRT:                   m.CRT,
		IsSimpleNational:      m.IsSimpleNational,
		MunicipalRegistration: m.MunicipalRegistration,
		ShowTaxBreakdown:      m.ShowTaxBreakdown,
		SendEmailToRecipient:  m.SendEmailToRecipient,
		BusinessName:          m.BusinessName,
		TradeName:             m.TradeName,
		Cnpj:                  m.Cnpj,
		Email:                 m.Email,
		Phone:                 m.Phone,
		Address: fiscalsettingsentity.FiscalAddress{
			Street:       m.Street,
			Number:       m.Number,
			Complement:   m.Complement,
			Neighborhood: m.Neighborhood,
			City:         m.City,
			UF:           m.UF,
			Cep:          m.Cep,
		},
	}
}

func (m *FiscalSettings) FromDomain(d *fiscalsettingsentity.FiscalSettings) {
	m.ID = d.ID
	m.CompanyID = d.CompanyID
	m.CompanyRegistryID = d.CompanyRegistryID
	m.IsActive = d.IsActive
	m.StateRegistration = d.StateRegistration
	m.TaxRegime = d.TaxRegime
	m.CNAE = d.CNAE
	m.CRT = d.CRT
	m.IsSimpleNational = d.IsSimpleNational
	m.MunicipalRegistration = d.MunicipalRegistration
	m.ShowTaxBreakdown = d.ShowTaxBreakdown
	m.SendEmailToRecipient = d.SendEmailToRecipient
	m.BusinessName = d.BusinessName
	m.TradeName = d.TradeName
	m.Cnpj = d.Cnpj
	m.Email = d.Email
	m.Phone = d.Phone
	m.Street = d.Address.Street
	m.Number = d.Address.Number
	m.Complement = d.Address.Complement
	m.Neighborhood = d.Address.Neighborhood
	m.City = d.Address.City
	m.UF = d.Address.UF
	m.Cep = d.Address.Cep
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = &d.UpdatedAt
}
