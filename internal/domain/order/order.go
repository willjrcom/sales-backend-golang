package orderentity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
)

type Order struct {
	entity.Entity
	Name        string
	OrderNumber int
	Items       []itementity.Items
	Delivery    *DeliveryOrder
	TableOrder  *TableOrder
	AttendantID uuid.UUID
	Status      StatusOrder
	Observation string
	PaymentOrder
}
