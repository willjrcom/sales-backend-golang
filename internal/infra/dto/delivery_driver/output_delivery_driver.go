package deliverydriverdto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type DeliveryDriverOutput struct {
	ID uuid.UUID `json:"id"`
	orderentity.DeliveryDriverCommonAttributes
}

func (s *DeliveryDriverOutput) FromModel(model *orderentity.DeliveryDriver) {
	s.ID = model.ID
	s.DeliveryDriverCommonAttributes = model.DeliveryDriverCommonAttributes
}
