package employeedto

import (
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type UpdateEmployeeInput struct {
}

func (r *UpdateEmployeeInput) validate() error {

	return nil
}

func (r *UpdateEmployeeInput) UpdateModel(employee *employeeentity.Employee) error {
	if err := r.validate(); err != nil {
		return err
	}

	return nil
}
