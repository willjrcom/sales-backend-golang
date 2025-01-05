package orderentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type DeliveryDriver struct {
	entity.Entity
	bun.BaseModel `bun:"table:delivery_drivers,alias:driver"`
	DeliveryDriverCommonAttributes
}

type DeliveryDriverCommonAttributes struct {
	EmployeeID      uuid.UUID                `bun:"column:employee_id,type:uuid,notnull" json:"employee_id"`
	Employee        *employeeentity.Employee `bun:"rel:belongs-to" json:"employee,omitempty"`
	OrderDeliveries []OrderDelivery          `bun:"rel:has-many,join:employee_id=driver_id" json:"order_deliveries,omitempty"`
}

type PatchDeliveryDriver struct {
}

func NewDeliveryDriver(deliveryDriverCommonAttributes DeliveryDriverCommonAttributes) *DeliveryDriver {
	return &DeliveryDriver{
		Entity:                         entity.NewEntity(),
		DeliveryDriverCommonAttributes: deliveryDriverCommonAttributes,
	}
}
