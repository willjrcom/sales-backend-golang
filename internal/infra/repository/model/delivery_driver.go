package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type DeliveryDriver struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:delivery_drivers,alias:driver"`
	DeliveryDriverCommonAttributes
}

type DeliveryDriverCommonAttributes struct {
	EmployeeID      uuid.UUID       `bun:"column:employee_id,type:uuid,notnull"`
	IsActive        bool            `bun:"column:is_active,type:boolean"`
	Employee        *Employee       `bun:"rel:belongs-to"`
	OrderDeliveries []OrderDelivery `bun:"rel:has-many,join:employee_id=driver_id"`
}

func (d *DeliveryDriver) FromDomain(driver *orderentity.DeliveryDriver) {
	if driver == nil {
		return
	}
	*d = DeliveryDriver{
		Entity: entitymodel.FromDomain(driver.Entity),
		DeliveryDriverCommonAttributes: DeliveryDriverCommonAttributes{
			EmployeeID: driver.EmployeeID,
			IsActive:   driver.IsActive,
			Employee:   &Employee{},
		},
	}

	d.Employee.FromDomain(driver.Employee)

	for _, orderDelivery := range driver.OrderDeliveries {
		od := OrderDelivery{}
		od.FromDomain(&orderDelivery)
		d.OrderDeliveries = append(d.OrderDeliveries, od)
	}
}

func (d *DeliveryDriver) ToDomain() *orderentity.DeliveryDriver {
	if d == nil {
		return nil
	}
	deliveryDriver := &orderentity.DeliveryDriver{
		Entity: d.Entity.ToDomain(),
		DeliveryDriverCommonAttributes: orderentity.DeliveryDriverCommonAttributes{
			EmployeeID:      d.EmployeeID,
			IsActive:        d.IsActive,
			Employee:        d.Employee.ToDomain(),
			OrderDeliveries: []orderentity.OrderDelivery{},
		},
	}

	for _, orderDelivery := range d.OrderDeliveries {
		deliveryDriver.OrderDeliveries = append(deliveryDriver.OrderDeliveries, *orderDelivery.ToDomain())
	}

	return deliveryDriver
}
