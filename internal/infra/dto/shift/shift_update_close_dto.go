package shiftdto

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrEndChangeRequired = errors.New("end change is required")
)

type ShiftUpdateCloseDTO struct {
	EndChange decimal.Decimal `json:"end_change"`
}

func (o *ShiftUpdateCloseDTO) validate() (err error) {
	if o.EndChange.IsZero() {
		return ErrEndChangeRequired
	}

	return
}

func (o *ShiftUpdateCloseDTO) ToDomain() (decimal.Decimal, error) {
	if err := o.validate(); err != nil {
		return decimal.Zero, err
	}

	return o.EndChange, nil
}
