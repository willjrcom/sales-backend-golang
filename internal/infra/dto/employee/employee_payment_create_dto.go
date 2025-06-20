package employeedto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

var (
	ErrAmountInvalid     = errors.New("amount must be non-negative")
	ErrMethodInvalid     = errors.New("payment method is invalid")
	ErrStatusInvalid     = errors.New("payment status is invalid")
	ErrEmployeeIDMissing = errors.New("employee_id is required")
)

type EmployeePaymentCreateDTO struct {
	EmployeeID  uuid.UUID                    `json:"employee_id"`
	Amount      decimal.Decimal              `json:"amount"`
	Status      employeeentity.PaymentStatus `json:"status"`
	Method      employeeentity.PaymentMethod `json:"method"`
	PaymentDate time.Time                    `json:"payment_date"`
	Notes       string                       `json:"notes"`
}

func (u *EmployeePaymentCreateDTO) validate() error {
	if u.Amount.IsNegative() {
		return ErrAmountInvalid
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
		u.PaymentDate,
		u.Notes,
	), nil
}
