package orderdto

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type CreateOrderInput struct {
	Name        string                     `json:"name"`
	Delivery    *orderentity.DeliveryOrder `json:"delivery"`
	TableOrder  *orderentity.TableOrder    `json:"table_order"`
	AttendantID uuid.UUID                  `json:"attendant_id"`
}

func (o *CreateOrderInput) Validate() error {
	if o.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func (o *CreateOrderInput) ToModel() (*orderentity.Order, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}

	return &orderentity.Order{
		Entity:      entity.Entity{ID: uuid.New(), CreatedAt: time.Now()},
		Name:        o.Name,
		Delivery:    o.Delivery,
		TableOrder:  o.TableOrder,
		AttendantID: o.AttendantID,
		Status:      orderentity.OrderStatusStaging,
	}, nil
}
