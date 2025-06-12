package orderdto

import (
	"errors"

	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrMethodInvalid    = errors.New("payment method is invalid")
	ErrTotalPaidInvalid = errors.New("total paid is invalid")
)

type OrderPaymentCreateDTO struct {
	TotalPaid decimal.Decimal       `json:"total_paid"`
	Method    orderentity.PayMethod `json:"method"`
}

func (u *OrderPaymentCreateDTO) validate() error {
	if u.TotalPaid.IsZero() || u.TotalPaid.IsNegative() {
		return ErrTotalPaidInvalid
	}

	for _, method := range orderentity.GetAllPayMethod() {
		if method == u.Method {
			return nil
		}

	}

	return ErrMethodInvalid
}

func (u *OrderPaymentCreateDTO) ToDomain(order *orderentity.Order) (*orderentity.PaymentOrder, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	return orderentity.NewPayment(u.TotalPaid, u.Method, order.ID), nil
}
