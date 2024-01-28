package orderentity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentOrder struct {
	PaymentTimeLogs
	TotalPaid float64   `bun:"total_paid" json:"total_paid"`
	Method    PayMethod `bun:"method,notnull" json:"method"`
	OrderID   uuid.UUID `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
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

func GetAllPayMethod() []PayMethod {
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
