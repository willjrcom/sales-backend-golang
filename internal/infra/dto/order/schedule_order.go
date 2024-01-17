package orderdto

import (
	"errors"
	"time"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrStartAtIsBeforeNow = errors.New("start at must be before now")
)

type UpdateScheduleOrder struct {
	orderentity.ScheduledOrder
}

func (o *UpdateScheduleOrder) validate() error {
	if o.StartAt != nil && o.StartAt.Before(time.Now()) {
		return ErrStartAtIsBeforeNow
	}

	return nil
}

func (o *UpdateScheduleOrder) ToModel() (*time.Time, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return o.StartAt, nil
}
