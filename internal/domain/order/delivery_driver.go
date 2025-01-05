package orderentity

import (
	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type DeliveryDriver struct {
	entity.Entity
	DeliveryDriverCommonAttributes
}

type DeliveryDriverCommonAttributes struct {
	EmployeeID      uuid.UUID
	Employee        *employeeentity.Employee
	OrderDeliveries []OrderDelivery
}

type PatchDeliveryDriver struct {
}

func NewDeliveryDriver(deliveryDriverCommonAttributes DeliveryDriverCommonAttributes) *DeliveryDriver {
	return &DeliveryDriver{
		Entity:                         entity.NewEntity(),
		DeliveryDriverCommonAttributes: deliveryDriverCommonAttributes,
	}
}
