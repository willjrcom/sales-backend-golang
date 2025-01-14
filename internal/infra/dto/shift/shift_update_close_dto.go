package shiftdto

import (
	"errors"
)

var (
	ErrEndChangeRequired = errors.New("end change is required")
)

type ShiftUpdateCloseDTO struct {
	EndChange *float64 `json:"end_change"`
}

func (o *ShiftUpdateCloseDTO) validate() (err error) {
	if o.EndChange == nil || *o.EndChange == 0 {
		return ErrEndChangeRequired
	}

	return
}

func (o *ShiftUpdateCloseDTO) ToDomain() (endChange float64, err error) {
	if err = o.validate(); err != nil {
		return 0, err
	}

	endChange = *o.EndChange
	return
}
