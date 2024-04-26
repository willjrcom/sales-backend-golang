package orderdto

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrMethodInvalid    = errors.New("payment method is invalid")
	ErrTotalPaidInvalid = errors.New("total paid is invalid")
)

type AddPaymentMethod struct {
	orderentity.PaymentOrder
}

func (u *AddPaymentMethod) validate() error {
	if err := u.validatePayMethod(); err != nil {
		return err
	}

	return nil
}

func (u *AddPaymentMethod) validatePayMethod() error {
	if u.TotalPaid <= 0 {
		return ErrTotalPaidInvalid
	}

	for _, method := range orderentity.GetAllPayMethod() {
		if method == u.Method {
			return nil
		}

	}

	return ErrMethodInvalid
}

func (u *AddPaymentMethod) ToModel(order *orderentity.Order) (*orderentity.PaymentOrder, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}

	return orderentity.NewPayment(u.TotalPaid, u.Method, order.ID), nil
}
