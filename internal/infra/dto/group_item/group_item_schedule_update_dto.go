package groupitemdto

import (
	"errors"
	"time"
)

var (
	ErrStartAtIsBeforeNow = errors.New("start at must be before now")
)

type GroupItemScheduleUpdateDTO struct {
	StartAt *time.Time `json:"start_at"`
}

func (o *GroupItemScheduleUpdateDTO) validate() error {
	if o.StartAt != nil && o.StartAt.Before(time.Now().UTC()) {
		return ErrStartAtIsBeforeNow
	}

	return nil
}

func (o *GroupItemScheduleUpdateDTO) ToDomain() (*time.Time, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return o.StartAt, nil
}
