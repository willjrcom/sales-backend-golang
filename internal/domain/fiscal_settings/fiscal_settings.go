package fiscalsettingsentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type FiscalSettings struct {
	entity.Entity
	CompanyID             uuid.UUID
	CompanyRegistryID     int64 // ID returned by Focus NFe
	TokenProduction       string
	TokenHomologation     string
	CSCProductionID       string
	CSCProductionCode     string
	CSCHomologationID     string
	CSCHomologationCode   string
	IsActive              bool
	StateRegistration     string // InscricaoEstadual
	TaxRegime             int    // RegimeTributario (1=Simples Nacional, 3=Regime Normal)
	CNAE                  string // CNAE
	CRT                   int    // CRT
	IsSimpleNational      bool   // SimplesNacional
	MunicipalRegistration string // InscricaoMunicipal

	// Preferences
	ShowTaxBreakdown     bool // DiscriminaImpostos
	SendEmailToRecipient bool // EnviarEmailDestinatario

	// Company Identity (Specific for Fiscal Emission)
	BusinessName string
	TradeName    string
	Cnpj         string
	Email        string
	Phone        string

	// Address (Specific for Fiscal Emission)
	Address FiscalAddress
}

type FiscalAddress struct {
	Street       string
	Number       string
	Complement   string
	Neighborhood string
	City         string
	UF           string
	Cep          string
}

func NewFiscalSettings(companyID uuid.UUID) *FiscalSettings {
	return &FiscalSettings{
		Entity:    entity.NewEntity(),
		CompanyID: companyID,
	}
}

func (f *FiscalSettings) Update(
	isActive bool,
	ie string,
	regime int,
	cnae string,
	crt int,
	businessName, tradeName, cnpj, email, phone string,
	address FiscalAddress,
	cscProductionID, cscProductionCode, cscHomologationID, cscHomologationCode string,
) {
	f.IsActive = isActive
	f.StateRegistration = ie
	f.TaxRegime = regime
	f.CNAE = cnae
	f.CRT = crt
	f.IsSimpleNational = regime == 1 || regime == 2
	f.BusinessName = businessName
	f.TradeName = tradeName
	f.Cnpj = cnpj
	f.Email = email
	f.Phone = phone
	f.Address = address
	f.CSCProductionID = cscProductionID
	f.CSCProductionCode = cscProductionCode
	f.CSCHomologationID = cscHomologationID
	f.CSCHomologationCode = cscHomologationCode
	f.UpdatedAt = time.Now().UTC()
}

func (f *FiscalSettings) SetCompanyRegistryID(id int64) {
	f.CompanyRegistryID = id
	f.UpdatedAt = time.Now().UTC()
}

func (f *FiscalSettings) SetTokens(production, homologation string) {
	f.TokenProduction = production
	f.TokenHomologation = homologation
	f.UpdatedAt = time.Now().UTC()
}
