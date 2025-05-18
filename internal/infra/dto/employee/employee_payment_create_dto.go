package employeedto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

var (
	ErrAmountInvalid     = errors.New("amount must be non-negative")
	ErrMethodInvalid     = errors.New("payment method is invalid")
	ErrStatusInvalid     = errors.New("payment status is invalid")
	ErrEmployeeIDMissing = errors.New("employee_id is required")
)

type EmployeePaymentCreateDTO struct {
	EmployeeID uuid.UUID                    `json:"employee_id"`
	Amount     float64                      `json:"amount"`
	Status     employeeentity.PaymentStatus `json:"status"`
	Method     employeeentity.PaymentMethod `json:"method"`
	PayDate    time.Time                    `json:"pay_date"`
	Notes      string                       `json:"notes"`
}

func (u *EmployeePaymentCreateDTO) validate() error {
	if u.Amount < 0 {
		return ErrAmountInvalid
	}
	validMethod := false
	for _, m := range employeeentity.GetAllPaymentMethods() {
		if m == u.Method {
			validMethod = true
			break
		}
	}
	if !validMethod {
		return ErrMethodInvalid
	}
	validStatus := false
	for _, s := range employeeentity.GetAllPaymentStatus() {
		if s == u.Status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return ErrStatusInvalid
	}
	if u.EmployeeID == uuid.Nil {
		return ErrEmployeeIDMissing
	}
	return nil
}

func (u *EmployeePaymentCreateDTO) ToDomain() (*employeeentity.PaymentEmployee, error) {
	if err := u.validate(); err != nil {
		return nil, err
	}
	return employeeentity.NewPaymentEmployee(
		u.EmployeeID,
		u.Amount,
		u.Status,
		u.Method,
		u.PayDate,
		u.Notes,
	), nil
}
