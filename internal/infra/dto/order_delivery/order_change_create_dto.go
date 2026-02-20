package orderdeliverydto

import (
	"errors"

	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrMethodInvalid = errors.New("payment method is invalid")
	ErrChangeInvalid = errors.New("change is invalid")
)

type OrderChangeCreateDTO struct {
	Change        decimal.Decimal       `json:"change"`
	PaymentMethod orderentity.PayMethod `json:"payment_method"`
}

func (u *OrderChangeCreateDTO) validate() error {
	if err := u.validatePayMethod(); err != nil {
		return err
	}

	return nil
}

func (u *OrderChangeCreateDTO) validatePayMethod() error {
	for _, method := range orderentity.GetAllPayMethod() {
		if method == u.PaymentMethod {
			return nil
		}

	}

	return ErrMethodInvalid
}

func (u *OrderChangeCreateDTO) ToDomain() (decimal.Decimal, *orderentity.PayMethod, error) {
	if err := u.validate(); err != nil {
		return decimal.Zero, nil, err
	}

	return u.Change, &u.PaymentMethod, nil
}
