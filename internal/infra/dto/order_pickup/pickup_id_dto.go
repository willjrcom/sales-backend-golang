package orderpickupdto

import "github.com/google/uuid"

type PickupIDAndOrderIDDTO struct {
	PickupID uuid.UUID `json:"pickup_id"`
	OrderID  uuid.UUID `json:"order_id"`
}

func FromDomain(pickupID uuid.UUID, orderID uuid.UUID) *PickupIDAndOrderIDDTO {
	return &PickupIDAndOrderIDDTO{
		PickupID: pickupID,
		OrderID:  orderID,
	}
}
