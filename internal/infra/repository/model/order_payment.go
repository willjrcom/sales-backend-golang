package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type PaymentOrder struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_payments,alias:payment"`
	PaymentTimeLogs
	PaymentCommonAttributes
}

type PaymentCommonAttributes struct {
	TotalPaid float64   `bun:"total_paid"`
	Method    string    `bun:"method,notnull"`
	OrderID   uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
}

type PaymentTimeLogs struct {
	PaidAt time.Time `bun:"paid_at"`
}
