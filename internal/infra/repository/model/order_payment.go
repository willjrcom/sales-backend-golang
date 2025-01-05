package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type PaymentOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:order_payments,alias:payment"`
	PaymentTimeLogs
	PaymentCommonAttributes
}

type PaymentCommonAttributes struct {
	TotalPaid float64   `bun:"total_paid"`
	Method    PayMethod `bun:"method,notnull"`
	OrderID   uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
}

type PaymentTimeLogs struct {
	PaidAt time.Time `bun:"paid_at"`
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
