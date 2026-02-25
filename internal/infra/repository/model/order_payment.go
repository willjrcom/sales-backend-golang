package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	TotalPaid *decimal.Decimal `bun:"total_paid,type:decimal(10,2)"`
	Method    string           `bun:"method,notnull"`
	OrderID   uuid.UUID        `bun:"column:order_id,type:uuid,notnull"`
}

type PaymentTimeLogs struct {
	PaidAt time.Time `bun:"paid_at"`
}

func (p *PaymentOrder) FromDomain(payment *orderentity.PaymentOrder) {
	if payment == nil {
		return
	}
	*p = PaymentOrder{
		Entity: entitymodel.FromDomain(payment.Entity),
		PaymentCommonAttributes: PaymentCommonAttributes{
			TotalPaid: &payment.TotalPaid,
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
			TotalPaid: p.GetTotalPaid(),
			Method:    orderentity.PayMethod(p.Method),
			OrderID:   p.OrderID,
		},
		PaymentTimeLogs: orderentity.PaymentTimeLogs{
			PaidAt: p.PaidAt,
		},
	}
}

func (p *PaymentOrder) GetTotalPaid() decimal.Decimal {
	if p.TotalPaid == nil {
		return decimal.Zero
	}
	return *p.TotalPaid
}
