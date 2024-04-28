package processdto

import (
	"errors"

	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
)

var (
	ErrIItemIDRequired       = errors.New("Item ID is required")
	ErrProcessRuleIDRequired = errors.New("process rule ID is required")
)

type CreateProcessInput struct {
	processentity.ProcessCommonAttributes
}

func (s *CreateProcessInput) validate() error {
	if s.ProcessRuleID == uuid.Nil {
		return ErrProcessRuleIDRequired
	}

	if s.ItemID == uuid.Nil {
		return ErrIItemIDRequired
	}

	return nil
}

func (s *CreateProcessInput) ToModel() (*processentity.Process, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	processCommonAttributes := processentity.ProcessCommonAttributes{
		ItemID:        s.ItemID,
		ProcessRuleID: s.ProcessRuleID,
	}

	return processentity.NewProcess(processCommonAttributes), nil
}
