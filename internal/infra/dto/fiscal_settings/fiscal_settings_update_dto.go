package fiscalsettingsdto

type FiscalSettingsUpdateDTO struct {
	FiscalEnabled      *bool   `json:"fiscal_enabled,omitempty"`
	InscricaoEstadual  *string `json:"inscricao_estadual,omitempty"`
	RegimeTributario   *int    `json:"regime_tributario,omitempty"`
	CNAE               *string `json:"cnae,omitempty"`
	CRT                *int    `json:"crt,omitempty"`
	InscricaoMunicipal *string `json:"inscricao_municipal,omitempty"`

	// Preferences
	DiscriminaImpostos      *bool `json:"discrimina_impostos,omitempty"`
	EnviarEmailDestinatario *bool `json:"enviar_email_destinatario,omitempty"`

	// Company Identity
	BusinessName *string `json:"business_name,omitempty"`
	TradeName    *string `json:"trade_name,omitempty"`
	Cnpj         *string `json:"cnpj,omitempty"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`

	// Address
	Street       *string `json:"street,omitempty"`
	Number       *string `json:"number,omitempty"`
	Complement   *string `json:"complement,omitempty"`
	Neighborhood *string `json:"neighborhood,omitempty"`
	City         *string `json:"city,omitempty"`
	UF           *string `json:"uf,omitempty"`
	Cep          *string `json:"cep,omitempty"`
}
