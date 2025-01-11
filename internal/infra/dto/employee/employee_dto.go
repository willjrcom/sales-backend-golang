package employeedto

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
)

type EmployeeDTO struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	companydto.UserDTO
}

func (c *EmployeeDTO) FromDomain(employee *employeeentity.Employee) {
	if employee == nil {
		return
	}

	*c = EmployeeDTO{
		ID:      employee.ID,
		UserID:  employee.UserID,
		UserDTO: companydto.UserDTO{},
	}

	c.UserDTO.FromDomain(employee.User)
}
