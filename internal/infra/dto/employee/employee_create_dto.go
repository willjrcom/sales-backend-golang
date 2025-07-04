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
	UserID      *uuid.UUID      `json:"user_id"`
	Permissions map[string]bool `json:"permissions"`
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

	employee := employeeentity.NewEmployee(*r.UserID)

	// Copia as permissões do DTO para o domain
	if r.Permissions != nil {
		employee.Permissions = make(employeeentity.Permissions)
		for k, v := range r.Permissions {
			employee.Permissions[employeeentity.PermissionKey(k)] = v
		}
	}

	return employee, nil
}
