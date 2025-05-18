package orderdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type PaymentOrderDTO struct {
	ID uuid.UUID `json:"id"`
	PaymentTimeLogs
	PaymentCommonAttributes
}

type PaymentCommonAttributes struct {
	TotalPaid decimal.Decimal       `json:"total_paid"`
	Method    orderentity.PayMethod `json:"method"`
	OrderID   uuid.UUID             `json:"order_id"`
}

type PaymentTimeLogs struct {
	PaidAt time.Time `json:"paid_at"`
}

func (p *PaymentOrderDTO) FromDomain(payment *orderentity.PaymentOrder) {
	if payment == nil {
		return
	}
	*p = PaymentOrderDTO{
		PaymentCommonAttributes: PaymentCommonAttributes{
			TotalPaid: payment.TotalPaid,
			Method:    payment.Method,
			OrderID:   payment.OrderID,
		},
		PaymentTimeLogs: PaymentTimeLogs{
			PaidAt: payment.PaidAt,
		},
		ID: payment.ID,
	}
}
