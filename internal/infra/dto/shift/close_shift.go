package shiftdto

import (
	"errors"

	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

var (
	ErrEndChangeRequired = errors.New("end change is required")
)

type CloseShift struct {
	shiftentity.ShiftCommonAttributes
}

func (o *CloseShift) validate() (err error) {
	if o.EndChange == nil || *o.EndChange == 0 {
		return ErrEndChangeRequired
	}

	return
}

func (o *CloseShift) ToModel() (endChange float32, err error) {
	if err = o.validate(); err != nil {
		return 0, err
	}

	endChange = *o.EndChange
	return
}
