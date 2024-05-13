package ordertabledto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type UpdateOrderTableInput struct {
	orderentity.OrderTableCommonAttributes
	ForceUpdate bool `json:"force_update"`
}

func (o *UpdateOrderTableInput) validate() error {
	return nil
}

func (s *UpdateOrderTableInput) ToModel() (model *orderentity.OrderTable, err error) {
	if err = s.validate(); err != nil {
		return nil, err
	}

	tableCommonAttributes := orderentity.OrderTableCommonAttributes{}

	if s.TableID != uuid.Nil {
		tableCommonAttributes.TableID = s.TableID
	}

	return &orderentity.OrderTable{
		OrderTableCommonAttributes: tableCommonAttributes,
	}, err
}
