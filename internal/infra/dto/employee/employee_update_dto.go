package employeedto

import (
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type EmployeeUpdateDTO struct {
}

func (r *EmployeeUpdateDTO) validate() error {
	return nil
}

func (r *EmployeeUpdateDTO) UpdateModel(employee *employeeentity.Employee) error {
	if err := r.validate(); err != nil {
		return err
	}

	return nil
}
