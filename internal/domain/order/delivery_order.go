package orderentity

import (
	"time"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type DeliveryOrder struct {
	Client      cliententity.Client
	Driver      *employeeentity.Employee
	Pickup      *time.Time
	Delivered   *time.Time
	DeliveryTax float64
	OrderID     uuid.UUID
}

func NewDeliveryOrder(Client *cliententity.Client, tax float64) *DeliveryOrder {
	return &DeliveryOrder{Client: *Client, DeliveryTax: tax}
}

func (d *DeliveryOrder) LaunchDelivery(driver *employeeentity.Employee) {
	*d.Driver = *driver
	*d.Pickup = time.Now()
}

func (d *DeliveryOrder) FinishDelivery() {
	*d.Delivered = time.Now()
}
