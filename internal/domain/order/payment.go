package orderentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type PaymentOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:payment_orders,alias:payment"`
	PaymentTimeLogs
	PaymentCommonAttributes
}

type PaymentCommonAttributes struct {
	TotalPaid float64   `bun:"total_paid" json:"total_paid"`
	Method    PayMethod `bun:"method,notnull" json:"method"`
	OrderID   uuid.UUID `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type PaymentTimeLogs struct {
	PaidAt time.Time `bun:"paid_at" json:"paid_at,omitempty"`
}

func NewPayment(totalPaid float64, method PayMethod, orderID uuid.UUID) *PaymentOrder {
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
