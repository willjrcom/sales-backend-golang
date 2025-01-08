package employeedto

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
)

type EmployeeDTO struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	userdto.UserDTO
}

func (c *EmployeeDTO) FromDomain(employee *employeeentity.Employee) {
	if employee == nil {
		return
	}

	*c = EmployeeDTO{
		ID:      employee.ID,
		UserID:  employee.UserID,
		UserDTO: userdto.UserDTO{},
	}

	c.UserDTO.FromDomain(employee.User)
}
