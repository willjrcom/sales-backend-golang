package fiscalsettingsdto

type FiscalSettingsUpdateDTO struct {
	FiscalEnabled         *bool   `json:"fiscal_enabled,omitempty"`
	StateRegistration     *string `json:"state_registration,omitempty"`
	TaxRegime             *int    `json:"tax_regime,omitempty"`
	CNAE                  *string `json:"cnae,omitempty"`
	CRT                   *int    `json:"crt,omitempty"`
	MunicipalRegistration *string `json:"municipal_registration,omitempty"`

	// Preferences
	ShowTaxBreakdown     *bool `json:"show_tax_breakdown,omitempty"`
	SendEmailToRecipient *bool `json:"send_email_to_recipient,omitempty"`

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
