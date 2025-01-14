package shiftdto

import (
	"errors"

	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

var (
	ErrNameRequired  = errors.New("name is required")
	ErrValueRequired = errors.New("value is required")
)

type ShiftRedeemCreateDTO struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (r *ShiftRedeemCreateDTO) validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Value == 0 {
		return ErrValueRequired
	}

	return nil
}

func (r *ShiftRedeemCreateDTO) ToDomain() (*shiftentity.Redeem, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	return &shiftentity.Redeem{
		Name:  r.Name,
		Value: r.Value,
	}, nil
}
