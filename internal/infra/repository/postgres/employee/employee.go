package employeerepositorylocal

import (
	"context"

	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type ProductRepositoryBun struct {
}

func NewProductRepositoryBun() *ProductRepositoryBun {
	return &ProductRepositoryBun{}
}

func (r *ProductRepositoryBun) RegisterEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return nil
}

func (r *ProductRepositoryBun) UpdateEmployee(ctx context.Context, p *employeeentity.Employee) error {
	return nil
}

func (r *ProductRepositoryBun) DeleteEmployee(ctx context.Context, id string) error {
	return nil
}

func (r *ProductRepositoryBun) GetEmployeeById(ctx context.Context, id string) (*employeeentity.Employee, error) {
	return nil, nil
}

func (r *ProductRepositoryBun) GetEmployeeBy(ctx context.Context, p *employeeentity.Employee) ([]employeeentity.Employee, error) {
	return nil, nil
}

func (r *ProductRepositoryBun) GetAllEmployee(ctx context.Context) ([]employeeentity.Employee, error) {
	return nil, nil
}
