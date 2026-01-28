package fiscalsettingsentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type FiscalSettings struct {
	entity.Entity
	CompanyID          uuid.UUID
	IsActive           bool
	InscricaoEstadual  string
	RegimeTributario   int
	CNAE               string
	CRT                int
	SimplesNacional    bool
	InscricaoMunicipal string

	// Preferences
	DiscriminaImpostos      bool
	EnviarEmailDestinatario bool

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
) {
	f.IsActive = isActive
	f.InscricaoEstadual = ie
	f.RegimeTributario = regime
	f.CNAE = cnae
	f.CRT = crt
	f.SimplesNacional = regime == 1 || regime == 2
	f.BusinessName = businessName
	f.TradeName = tradeName
	f.Cnpj = cnpj
	f.Email = email
	f.Phone = phone
	f.Address = address
	f.UpdatedAt = time.Now()
}
