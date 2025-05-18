package model

import "context"

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, p *Employee) error
	UpdateEmployee(ctx context.Context, p *Employee) error
	DeleteEmployee(ctx context.Context, id string) error
	GetEmployeeById(ctx context.Context, id string) (*Employee, error)
	GetEmployeeByUserID(ctx context.Context, userID string) (*Employee, error)
	GetAllEmployees(ctx context.Context) ([]Employee, error)
	// AddPaymentEmployee records a payment for an employee.
	AddPaymentEmployee(ctx context.Context, p *PaymentEmployee) error
}
