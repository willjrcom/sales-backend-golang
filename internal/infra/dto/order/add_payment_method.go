package orderdto

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrMethodInvalid = errors.New("payment method is invalid")
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

	u.OrderID = order.ID
	return &u.PaymentOrder, nil
}