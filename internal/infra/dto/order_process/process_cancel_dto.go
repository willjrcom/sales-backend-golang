package processdto

import (
	"errors"
)

var (
	ErrReasonIsRequired = errors.New("reason is required")
)

type OrderProcessCancelDTO struct {
	Reason *string `json:"reason"`
}

func (s *OrderProcessCancelDTO) validate() error {
	if s.Reason == nil {
		return ErrReasonIsRequired
	}

	return nil
}

func (s *OrderProcessCancelDTO) ToDomain() (*string, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	return s.Reason, nil
}
