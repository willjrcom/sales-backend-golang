package employeedto

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type EmployeeOutput struct {
	ID uuid.UUID `json:"id"`
	employeeentity.Employee
}

func (c *EmployeeOutput) FromModel(model *employeeentity.Employee) {
	c.ID = model.ID
	c.Employee = *model
}
