package employeerepositorylocal

import (
	"context"

	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type ProductRepositoryLocal struct {
}

func NewProductRepositoryLocal() *ProductRepositoryLocal {
	return &ProductRepositoryLocal{}
}

func (r *ProductRepositoryLocal) RegisterEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return nil
}

func (r *ProductRepositoryLocal) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return nil
}

func (r *ProductRepositoryLocal) DeleteEmployee(ctx context.Context, id string) error {
	return nil
}

func (r *ProductRepositoryLocal) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetEmployeeBy(ctx context.Context, p *employeeentity.Employee) (*employeeentity.Employee, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetAllEmployee(ctx context.Context) ([]employeeentity.Employee, error) {
	return nil, nil
}
