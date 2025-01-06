package deliverydriverdto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
)

type DeliveryDriverDTO struct {
	ID              uuid.UUID                   `json:"id"`
	EmployeeID      uuid.UUID                   `json:"employee_id"`
	Employee        *employeedto.EmployeeDTO    `json:"employee"`
	OrderDeliveries []orderentity.OrderDelivery `json:"order_deliveries"`
}

func (s *DeliveryDriverDTO) FromModel(model *orderentity.DeliveryDriver) {
	*s = DeliveryDriverDTO{
		ID:              model.ID,
		EmployeeID:      model.EmployeeID,
		OrderDeliveries: model.OrderDeliveries,
	}

	s.Employee.FromModel(model.Employee)
}
