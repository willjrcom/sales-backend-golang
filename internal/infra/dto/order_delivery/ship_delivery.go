package orderdeliverydto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type ShipDeliveryOrder struct {
	DriverID    uuid.UUID `json:"driver_id"`
	DeliveryIDs []string  `json:"delivery_ids"`
}

func (u *ShipDeliveryOrder) validate() error {
	if u.DriverID == uuid.Nil {
		return ErrInvalidDriverID
	}
	return nil
}

func (u *ShipDeliveryOrder) UpdateModel(deliveries []orderentity.OrderDelivery) error {
	if err := u.validate(); err != nil {
		return err
	}

	for i := range deliveries {
		deliveries[i].DriverID = &u.DriverID
	}

	return nil
}
