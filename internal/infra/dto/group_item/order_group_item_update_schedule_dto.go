package groupitemdto

import (
	"errors"
	"time"
)

var (
	ErrStartAtIsBeforeNow = errors.New("start at must be before now")
)

type OrderGroupItemUpdateScheduleDTO struct {
	StartAt *time.Time `json:"start_at"`
}

func (o *OrderGroupItemUpdateScheduleDTO) validate() error {
	if o.StartAt != nil && o.StartAt.Before(time.Now().UTC()) {
		return ErrStartAtIsBeforeNow
	}

	return nil
}

func (o *OrderGroupItemUpdateScheduleDTO) ToDomain() (*time.Time, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return o.StartAt, nil
}
