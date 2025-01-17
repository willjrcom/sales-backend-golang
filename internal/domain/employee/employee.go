package employeeentity

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Employee struct {
	entity.Entity
	UserID uuid.UUID
	User   *companyentity.User
}

func NewEmployee(userID uuid.UUID) *Employee {
	return &Employee{
		Entity: entity.NewEntity(),
		UserID: userID,
	}
}
