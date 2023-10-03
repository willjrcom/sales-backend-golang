package orderentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type DeliveryOrder struct {
	entity.Entity
	bun.BaseModel `bun:"table:delivery_orders"`
	OrderID       uuid.UUID                `bun:"column:order_id,type:uuid"`
	Pickup        *time.Time               `bun:"pickup"`
	Delivered     *time.Time               `bun:"delivered"`
	IsCompleted   bool                     `bun:"is_completed"`
	DeliveryTax   *float64                 `bun:"delivery_tax"`
	ClientID      uuid.UUID                `bun:"column:client_id,type:uuid,notnull"`
	Client        *cliententity.Client     `bun:"rel:belongs-to"`
	AddressID     uuid.UUID                `bun:"column:address_id,type:uuid,notnull"`
	Address       *addressentity.Address   `bun:"rel:belongs-to"`
	DriverID      *uuid.UUID               `bun:"column:driver_id,type:uuid,notnull"`
	Driver        *employeeentity.Employee `bun:"rel:belongs-to"`
}

func (d *DeliveryOrder) LaunchDelivery(driver *employeeentity.Employee) {
	*d.Driver = *driver
	*d.Pickup = time.Now()
}

func (d *DeliveryOrder) FinishDelivery() {
	*d.Delivered = time.Now()
}
