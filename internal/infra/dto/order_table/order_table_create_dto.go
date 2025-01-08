package ordertabledto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrOrderIDRequired = errors.New("order id is required")
	ErrTableIDRequired = errors.New("table_id is required")
)

type CreateOrderTableInput struct {
	Name    string    `json:"name,omitempty"`
	Contact string    `json:"contact,omitempty"`
	TableID uuid.UUID `json:"table_id"`
}

func (o *CreateOrderTableInput) validate() error {
	if o.TableID == uuid.Nil {
		return ErrTableIDRequired
	}

	return nil
}

func (o *CreateOrderTableInput) ToDomain() (*orderentity.OrderTable, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	orderTableCommonAttributes := orderentity.OrderTableCommonAttributes{
		Name:    o.Name,
		Contact: o.Contact,
		TableID: o.TableID,
	}
	return orderentity.NewTable(orderTableCommonAttributes), nil
}
