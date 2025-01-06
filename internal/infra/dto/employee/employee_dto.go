package employeedto

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
)

type EmployeeOutput struct {
	ID uuid.UUID `json:"id"`
	userdto.UserDTO
}

func (c *EmployeeOutput) FromModel(employee *employeeentity.Employee) {
	c.UserDTO.FromModel(employee.User)
}
