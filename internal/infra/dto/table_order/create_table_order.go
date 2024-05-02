package tableorderdto

import (
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrOrderIDRequired = errors.New("order id is required")
	ErrTableIDRequired = errors.New("table_id is required")
)

type CreateTableOrderInput struct {
	orderentity.TableOrderCommonAttributes
}

func (o *CreateTableOrderInput) validate() error {
	if o.OrderID == uuid.Nil {
		return ErrOrderIDRequired
	}

	if o.TableID == uuid.Nil {
		return ErrTableIDRequired
	}

	return nil
}

func (o *CreateTableOrderInput) ToModel() (*orderentity.TableOrder, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return orderentity.NewTable(o.TableOrderCommonAttributes), nil
}
