package orderentity

import "time"

type PaymentOrder struct {
	PaymentTimeLogs
	TotalPaid float64   `bun:"total_paid" json:"total_paid"`
	Change    float64   `bun:"change" json:"change"`
	Method    PayMethod `bun:"method,notnull" json:"method"`
}

type PatchPaymentOrder struct {
	PaymentTimeLogs
	TotalPaid *float64   `json:"total_paid"`
	Change    *float64   `json:"change"`
	Method    *PayMethod `json:"method"`
}

type PaymentTimeLogs struct {
	PaidAt *time.Time `bun:"paid_at" json:"paid_at,omitempty"`
}

type PayMethod string

// Tipos de cartão
const (
	Dinheiro        PayMethod = "Dinheiro"
	Visa            PayMethod = "Visa"
	MasterCard      PayMethod = "MasterCard"
	Ticket          PayMethod = "Ticket"
	VR              PayMethod = "VR"
	AmericanExpress PayMethod = "American Express"
	Elo             PayMethod = "Elo"
	DinersClub      PayMethod = "Diners Club"
	Hipercard       PayMethod = "Hipercard"
	VisaElectron    PayMethod = "Visa Electron"
	Maestro         PayMethod = "Maestro"
	Alelo           PayMethod = "Alelo"
	PayPal          PayMethod = "PayPal"
	Outros          PayMethod = "Outros"
)

func GetAll() []PayMethod {
	return []PayMethod{
		Visa,
		MasterCard,
		Ticket,
		VR,
		AmericanExpress,
		Elo,
		DinersClub,
		Hipercard,
		VisaElectron,
		Maestro,
		Alelo,
		PayPal,
		Outros,
	}
}
