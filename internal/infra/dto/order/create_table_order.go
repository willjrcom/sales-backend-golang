package orderdto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrCodTableRequired = errors.New("num_table is required")
	ErrWaiterIDRequired = errors.New("waiter_id is required")
)

type CreateTableOrderInput struct {
	OrderID  uuid.UUID `json:"order_id"`
	CodTable string    `json:"cod_table"`
	QrCode   *string   `json:"qr_code"`
	WaiterID uuid.UUID `json:"waiter_id"`
}

func (o *CreateTableOrderInput) Validate() error {
	if o.OrderID == uuid.Nil {
		return ErrOrderIDRequired
	}

	if o.CodTable == "" {
		return ErrCodTableRequired
	}

	if o.WaiterID == uuid.Nil {
		return ErrWaiterIDRequired
	}
	return nil
}

func (o *CreateTableOrderInput) ToModel() (*orderentity.TableOrder, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}

	table := &orderentity.TableOrder{
		Entity:   entity.NewEntity(),
		OrderID:  o.OrderID,
		CodTable: o.CodTable,
		WaiterID: o.WaiterID,
	}

	if o.QrCode != nil {
		table.QrCode = *o.QrCode
	}

	return table, nil
}
