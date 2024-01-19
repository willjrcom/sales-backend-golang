package orderdto

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrShiftIDRequired = errors.New("shift ID is required")
)

type CreateOrderInput struct {
	ShiftID *uuid.UUID `json:"shift_id"`
}

func (o *CreateOrderInput) validate() error {
	if o.ShiftID == nil {
		return ErrShiftIDRequired
	}

	return nil
}

func (o *CreateOrderInput) ToModel() (*uuid.UUID, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return o.ShiftID, nil
}
