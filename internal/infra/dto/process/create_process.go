package processdto

import (
	"errors"

	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
)

var (
	ErrItemIDRequired        = errors.New("item ID is required")
	ErrProcessRuleIDRequired = errors.New("process rule ID is required")
)

type CreateProcessInput struct {
	processentity.ProcessCommonAttributes
}

func (s *CreateProcessInput) validate() error {
	if s.ProcessRuleID == uuid.Nil {
		return ErrProcessRuleIDRequired
	}

	if s.GroupItemID == uuid.Nil {
		return ErrItemIDRequired
	}

	return nil
}

func (s *CreateProcessInput) ToModel() (*processentity.Process, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	return processentity.NewProcess(s.GroupItemID, s.ProcessRuleID), nil
}
