package employeedto

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
)

type EmployeeDTO struct {
	companydto.UserDTO
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	Permissions map[string]bool `json:"permissions"`
	IsActive    bool            `json:"is_active"`
}

func (c *EmployeeDTO) FromDomain(employee *employeeentity.Employee) {
	if employee == nil {
		return
	}

	*c = EmployeeDTO{
		ID:          employee.ID,
		UserID:      employee.UserID,
		UserDTO:     companydto.UserDTO{},
		Permissions: make(map[string]bool),
		IsActive:    employee.IsActive,
	}

	c.UserDTO.FromDomain(employee.User)

	// Copia as permissões do domain para o DTO
	for k, v := range employee.Permissions {
		c.Permissions[string(k)] = v
	}
}
