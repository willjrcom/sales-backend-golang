package employeedto

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
)

type EmployeeDTO struct {
	ID uuid.UUID `json:"id"`
	userdto.UserDTO
}

func (c *EmployeeDTO) FromModel(employee *employeeentity.Employee) {
	c.UserDTO.FromModel(employee.User)
}
