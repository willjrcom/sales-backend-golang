package employeeentity

import (
	"github.com/google/uuid"
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
