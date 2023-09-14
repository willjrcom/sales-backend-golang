package employeeentity

import "context"

type Repository interface {
	RegisterEmployee(ctx context.Context, p *Employee) error
	UpdateEmployee(ctx context.Context, p *Employee) error
	DeleteEmployee(ctx context.Context, id string) error
	GetEmployeeById(ctx context.Context, id string) (*Employee, error)
	GetEmployeeBy(ctx context.Context, p *Employee) ([]Employee, error)
	GetAllEmployee(ctx context.Context) ([]Employee, error)
}
