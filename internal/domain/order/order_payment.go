package orderentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type PaymentOrder struct {
	entity.Entity
	PaymentTimeLogs
	PaymentCommonAttributes
}

type PaymentCommonAttributes struct {
	TotalPaid decimal.Decimal
	Method    PayMethod
	OrderID   uuid.UUID
}

type PaymentTimeLogs struct {
	PaidAt time.Time
}

// NewPayment creates a new payment record for an order
func NewPayment(totalPaid decimal.Decimal, method PayMethod, orderID uuid.UUID) *PaymentOrder {
	return &PaymentOrder{
		Entity: entity.NewEntity(),
		PaymentCommonAttributes: PaymentCommonAttributes{
			TotalPaid: totalPaid,
			Method:    method,
			OrderID:   orderID,
		},
		PaymentTimeLogs: PaymentTimeLogs{
			PaidAt: time.Now().UTC(),
		},
	}
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
		Dinheiro,
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
