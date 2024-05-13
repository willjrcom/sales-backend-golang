package processdto

import (
	"errors"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

var (
	ErrGroupItemIDRequired   = errors.New("group item ID is required")
	ErrProcessRuleIDRequired = errors.New("process rule ID is required")
)

type CreateProcessInput struct {
	orderprocessentity.ProcessCommonAttributes
}

func (s *CreateProcessInput) validate() error {
	if s.ProcessRuleID == uuid.Nil {
		return ErrProcessRuleIDRequired
	}

	if s.GroupItemID == uuid.Nil {
		return ErrGroupItemIDRequired
	}

	return nil
}

func (s *CreateProcessInput) ToModel() (*orderprocessentity.Process, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	return orderprocessentity.NewProcess(s.GroupItemID, s.ProcessRuleID), nil
}
