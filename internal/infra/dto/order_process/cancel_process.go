package processdto

import (
	"errors"
)

var (
	ErrReasonIsRequired = errors.New("reason is required")
)

type CancelProcess struct {
	Reason *string `json:"reason"`
}

func (s *CancelProcess) validate() error {
	if s.Reason == nil {
		return ErrReasonIsRequired
	}

	return nil
}

func (s *CancelProcess) ToModel() (*string, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	return s.Reason, nil
}
