package employeerepositorylocal

import (
	"context"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type EmployeeRepositoryLocal struct {
}

func NewEmployeeRepositoryLocal() *EmployeeRepositoryLocal {
	return &EmployeeRepositoryLocal{}
}

func (r *EmployeeRepositoryLocal) CreateEmployee(ctx context.Context, p *model.Employee) error {
	return nil
}

func (r *EmployeeRepositoryLocal) UpdateEmployee(ctx context.Context, p *model.Employee) error {
	return nil
}

func (r *EmployeeRepositoryLocal) DeleteEmployee(ctx context.Context, id string) error {
	return nil
}

func (r *EmployeeRepositoryLocal) GetEmployeeById(ctx context.Context, id string) (*model.Employee, error) {
	return nil, nil
}

func (r *EmployeeRepositoryLocal) GetAllEmployees(ctx context.Context) ([]model.Employee, error) {
	return nil, nil
}
