package employeedto

import (
	"errors"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

var (
	ErrUserIDRequired = errors.New("user_id is required")
)

type EmployeeCreateDTO struct {
	UserID *uuid.UUID `json:"user_id"`
}

func (r *EmployeeCreateDTO) validate() error {
	if r.UserID == nil {
		return ErrUserIDRequired
	}

	return nil
}

func (r *EmployeeCreateDTO) ToDomain() (*employeeentity.Employee, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	// Get exists user
	return employeeentity.NewEmployee(*r.UserID), nil
}
