package orderentity

import (
	"time"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type DeliveryOrder struct {
	ClientID    uuid.UUID
	Client      *cliententity.Client `bun:"rel:belongs-to,join:client_id=id"`
	DriverID    uuid.UUID
	Driver      *employeeentity.Employee `bun:"rel:has-one,join:driver_id=id"`
	Pickup      *time.Time
	Delivered   *time.Time
	IsCompleted bool
	DeliveryTax float64
}

func (d *DeliveryOrder) LaunchDelivery(driver *employeeentity.Employee) {
	*d.Driver = *driver
	*d.Pickup = time.Now()
}

func (d *DeliveryOrder) FinishDelivery() {
	*d.Delivered = time.Now()
}
