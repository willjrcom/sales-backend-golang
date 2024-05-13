package processdto

import (
	"errors"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

var (
	ErrEmployeeIDRequired = errors.New("employee ID is required")
)

type StartProcessInput struct {
	orderprocessentity.OrderProcessCommonAttributes
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
