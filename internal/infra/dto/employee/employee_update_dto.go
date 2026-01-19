package employeedto

import (
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type EmployeeUpdateDTO struct {
	Permissions map[string]bool `json:"permissions"`
	IsActive    *bool           `json:"is_active"`
}

func (r *EmployeeUpdateDTO) validate() error {
	return nil
}

func (r *EmployeeUpdateDTO) UpdateDomain(employee *employeeentity.Employee) error {
	if err := r.validate(); err != nil {
		return err
	}

	// Atualiza as permissões se fornecidas
	if r.Permissions != nil {
		employee.Permissions = make(employeeentity.Permissions)
		for k, v := range r.Permissions {
			employee.Permissions[employeeentity.PermissionKey(k)] = v
		}
	}

	if r.IsActive != nil {
		employee.IsActive = *r.IsActive
	}
	return nil
}
