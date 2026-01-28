package fiscalsettingsdto

import (
	fiscalsettingsentity "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_settings"
)

type FiscalSettingsDTO struct {
	FiscalEnabled      bool   `json:"fiscal_enabled"`
	InscricaoEstadual  string `json:"inscricao_estadual"`
	RegimeTributario   int    `json:"regime_tributario"`
	CNAE               string `json:"cnae"`
	CRT                int    `json:"crt"`
	SimplesNacional    bool   `json:"simples_nacional"`
	InscricaoMunicipal string `json:"inscricao_municipal"`

	// Preferences
	DiscriminaImpostos      bool `json:"discrimina_impostos"`
	EnviarEmailDestinatario bool `json:"enviar_email_destinatario"`

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
	d.FiscalEnabled = entity.IsActive
	d.InscricaoEstadual = entity.InscricaoEstadual
	d.RegimeTributario = entity.RegimeTributario
	d.CNAE = entity.CNAE
	d.CRT = entity.CRT
	d.SimplesNacional = entity.SimplesNacional
	d.InscricaoMunicipal = entity.InscricaoMunicipal
	d.DiscriminaImpostos = entity.DiscriminaImpostos
	d.EnviarEmailDestinatario = entity.EnviarEmailDestinatario

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
