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

func (r *EmployeeRepositoryLocal) RegisterEmployee(ctx context.Context, p *employeeentity.Employee) error {
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

func (r *EmployeeRepositoryLocal) GetEmployeeBy(ctx context.Context, p *employeeentity.Employee) ([]employeeentity.Employee, error) {
	return nil, nil
}

func (r *EmployeeRepositoryLocal) GetAllEmployee(ctx context.Context) ([]employeeentity.Employee, error) {
	return nil, nil
}
