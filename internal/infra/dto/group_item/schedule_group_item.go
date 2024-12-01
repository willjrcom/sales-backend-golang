package groupitemdto

import (
	"errors"
	"time"
)

var (
	ErrStartAtIsBeforeNow = errors.New("start at must be before now")
)

type UpdateScheduleGroupItem struct {
	StartAt *time.Time `json:"start_at"`
}

func (o *UpdateScheduleGroupItem) validate() error {
	if o.StartAt != nil && o.StartAt.Before(time.Now()) {
		return ErrStartAtIsBeforeNow
	}

	return nil
}

func (o *UpdateScheduleGroupItem) ToModel() (*time.Time, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return o.StartAt, nil
}
