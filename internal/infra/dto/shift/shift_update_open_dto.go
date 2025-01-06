package shiftdto

import (
	"errors"
)

var (
	ErrStartChangeRequired = errors.New("start change must be higher than 0")
)

type ShiftUpdateOpenDTO struct {
	StartChange *float32 `json:"start_change"`
}

func (o *ShiftUpdateOpenDTO) validate() (err error) {
	if o.StartChange == nil || *o.StartChange == 0 {
		return ErrStartChangeRequired
	}

	return
}

func (o *ShiftUpdateOpenDTO) ToDomain() (startChange float32, err error) {
	if err = o.validate(); err != nil {
		return 0, err
	}

	return startChange, nil
}
