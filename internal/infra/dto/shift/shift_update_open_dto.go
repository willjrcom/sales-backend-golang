package shiftdto

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrStartChangeRequired = errors.New("start change must be higher than 0")
)

type ShiftUpdateOpenDTO struct {
	StartChange decimal.Decimal `json:"start_change"`
}

func (o *ShiftUpdateOpenDTO) validate() (err error) {
	if o.StartChange.IsZero() {
		return ErrStartChangeRequired
	}

	return
}

func (o *ShiftUpdateOpenDTO) ToDomain() (decimal.Decimal, error) {
	if err := o.validate(); err != nil {
		return decimal.Zero, err
	}

	return o.StartChange, nil
}
