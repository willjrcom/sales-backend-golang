package employeedto

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type EmployeeOutput struct {
	ID uuid.UUID `json:"id"`
	personentity.PersonCommonAttributes
}

func (c *EmployeeOutput) FromModel(model *employeeentity.Employee) {
	c.ID = model.ID
	c.PersonCommonAttributes = model.User.Person.PersonCommonAttributes
}
