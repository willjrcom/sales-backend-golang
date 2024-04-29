package processdto

import (
	"errors"

	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
)

var (
	ErrEmployeeIDRequired = errors.New("employee ID is required")
)

type StartProcessInput struct {
	processentity.ProcessCommonAttributes
}

func (s *StartProcessInput) validate() error {
	if s.EmployeeID == nil {
		return ErrEmployeeIDRequired
	}

	return nil
}

func (s *StartProcessInput) ToModel() (uuid.UUID, error) {
	if err := s.validate(); err != nil {
		return uuid.Nil, err
	}

	return *s.EmployeeID, nil
}
