package deliverydriverdto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
)

type DeliveryDriverDTO struct {
	ID              uuid.UUID                   `json:"id"`
	EmployeeID      uuid.UUID                   `json:"employee_id"`
	IsActive        bool                        `json:"is_active"`
	Employee        *employeedto.EmployeeDTO    `json:"employee"`
	OrderDeliveries []orderentity.OrderDelivery `json:"order_deliveries"`
}

func (s *DeliveryDriverDTO) FromDomain(driver *orderentity.DeliveryDriver) {
	if driver == nil {
		return
	}
	*s = DeliveryDriverDTO{
		ID:              driver.ID,
		EmployeeID:      driver.EmployeeID,
		IsActive:        driver.IsActive,
		Employee:        &employeedto.EmployeeDTO{},
		OrderDeliveries: driver.OrderDeliveries,
	}

	s.Employee.FromDomain(driver.Employee)

	if driver.Employee == nil {
		s.Employee = nil
	}
}
