package shiftdto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

var (
	ErrAttendantIDRequired         = errors.New("attendant id is required")
	ErrDayIsRequired               = errors.New("day is required")
	ErrStartChangeRequired         = errors.New("start change is required")
	ErrEndChangeNotUsedToOpenShift = errors.New("end change is not used to open shift")
)

type OpenShift struct {
	shiftentity.ShiftCommonAttributes
}

func (o *OpenShift) validate() (err error) {
	if o.AttendantID == nil {
		return ErrAttendantIDRequired
	}

	if o.Day == nil {
		return ErrDayIsRequired
	}

	if o.StartChange == 0 {
		return ErrStartChangeRequired
	}

	if o.EndChange != nil {
		return ErrEndChangeNotUsedToOpenShift
	}
	return
}

func (o *OpenShift) ToModel() (shift *shiftentity.Shift, err error) {
	if err = o.validate(); err != nil {
		return nil, err
	}
	return &shiftentity.Shift{
		Entity:                entity.NewEntity(),
		ShiftCommonAttributes: o.ShiftCommonAttributes,
	}, nil
}
