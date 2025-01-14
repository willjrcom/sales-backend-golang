package shiftdto

import (
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

type RedeemDTO struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (r *RedeemDTO) FromDomain(redeem *shiftentity.Redeem) {
	if redeem == nil {
		return
	}
	*r = RedeemDTO{
		Name:  redeem.Name,
		Value: redeem.Value,
	}
}
