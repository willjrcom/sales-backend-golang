package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type DeliveryDriver struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:delivery_drivers,alias:driver"`
	DeliveryDriverCommonAttributes
}

type DeliveryDriverCommonAttributes struct {
	EmployeeID      uuid.UUID       `bun:"column:employee_id,type:uuid,notnull"`
	Employee        *Employee       `bun:"rel:belongs-to"`
	OrderDeliveries []OrderDelivery `bun:"rel:has-many,join:employee_id=driver_id"`
}
