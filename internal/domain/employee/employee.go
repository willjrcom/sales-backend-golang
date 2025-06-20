package employeeentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Employee struct {
	entity.Entity
	UserID   uuid.UUID
	User     *companyentity.User
	Payments []PaymentEmployee
}

func NewEmployee(userID uuid.UUID) *Employee {
	return &Employee{
		Entity:   entity.NewEntity(),
		UserID:   userID,
		Payments: make([]PaymentEmployee, 0),
	}
}

type EmployeeSalaryHistory struct {
	entity.Entity
	EmployeeID uuid.UUID
	StartDate  time.Time
	EndDate    *time.Time
	SalaryType string
	BaseSalary decimal.Decimal
	HourlyRate decimal.Decimal
	Commission float64
}
