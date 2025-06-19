package shiftentity

import "github.com/shopspring/decimal"

type Redeem struct {
	Name  string
	Value decimal.Decimal
}

func NewRedeem(name string, value decimal.Decimal) *Redeem {
	return &Redeem{
		Name:  name,
		Value: value,
	}
}
