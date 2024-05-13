package employeerepositorylocal

import (
	"context"

	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type EmployeeRepositoryLocal struct {
}

func NewEmployeeRepositoryLocal() *EmployeeRepositoryLocal {
	return &EmployeeRepositoryLocal{}
}

func (r *EmployeeRepositoryLocal) CreateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return nil
}

func (r *EmployeeRepositoryLocal) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return nil
}

func (r *EmployeeRepositoryLocal) DeleteEmployee(ctx context.Context, id string) error {
	return nil
}

func (r *EmployeeRepositoryLocal) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	return nil, nil
}

func (r *EmployeeRepositoryLocal) GetAllEmployees(ctx context.Context) ([]employeeentity.Employee, error) {
	return nil, nil
}
