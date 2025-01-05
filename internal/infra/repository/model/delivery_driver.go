package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type DeliveryDriver struct {
	entity.Entity
	bun.BaseModel `bun:"table:delivery_drivers,alias:driver"`
	DeliveryDriverCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type DeliveryDriverCommonAttributes struct {
	EmployeeID      uuid.UUID                `bun:"column:employee_id,type:uuid,notnull"`
	Employee        *employeeentity.Employee `bun:"rel:belongs-to"`
	OrderDeliveries []OrderDelivery          `bun:"rel:has-many,join:employee_id=driver_id"`
}
