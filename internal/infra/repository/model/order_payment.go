package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
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

func (p *PaymentOrder) FromDomain(payment *orderentity.PaymentOrder) {
	*p = PaymentOrder{
		Entity: entitymodel.FromDomain(payment.Entity),
		PaymentCommonAttributes: PaymentCommonAttributes{
			TotalPaid: payment.TotalPaid,
			Method:    string(payment.Method),
			OrderID:   payment.OrderID,
		},
		PaymentTimeLogs: PaymentTimeLogs{
			PaidAt: payment.PaidAt,
		},
	}
}

func (p *PaymentOrder) ToDomain() *orderentity.PaymentOrder {
	if p == nil {
		return nil
	}
	return &orderentity.PaymentOrder{
		Entity: p.Entity.ToDomain(),
		PaymentCommonAttributes: orderentity.PaymentCommonAttributes{
			TotalPaid: p.TotalPaid,
			Method:    orderentity.PayMethod(p.Method),
			OrderID:   p.OrderID,
		},
		PaymentTimeLogs: orderentity.PaymentTimeLogs{
			PaidAt: p.PaidAt,
		},
	}
}
