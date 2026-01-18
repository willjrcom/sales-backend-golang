package model

import (
	"context"

	"github.com/google/uuid"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, p *Employee) error
	UpdateEmployee(ctx context.Context, p *Employee) error
	DeleteEmployee(ctx context.Context, id string) error
	GetEmployeeById(ctx context.Context, id string) (*Employee, error)
	GetEmployeeByUserID(ctx context.Context, userID string) (*Employee, error)
	GetEmployeeDeletedByUserID(ctx context.Context, userID string) (*Employee, error)
	GetAllEmployees(ctx context.Context, page, perPage int, isActive ...bool) ([]Employee, int, error)
	AddPaymentEmployee(ctx context.Context, p *PaymentEmployee) error
	GetAllEmployeeDeleted(ctx context.Context, page, perPage int) ([]Employee, int, error)
	GetSalaryHistory(ctx context.Context, employeeID uuid.UUID) ([]EmployeeSalaryHistory, error)
	GetPayments(ctx context.Context, employeeID uuid.UUID) ([]PaymentEmployee, error)
	CreateSalaryHistory(ctx context.Context, h *EmployeeSalaryHistory) error
}
