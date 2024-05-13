package orderpickupdto

import "github.com/google/uuid"

type PickupIDAndOrderIDOutput struct {
	PickupID uuid.UUID `json:"pickup_id"`
	OrderID  uuid.UUID `json:"order_id"`
}

func NewOutput(pickupID uuid.UUID, orderID uuid.UUID) *PickupIDAndOrderIDOutput {
	return &PickupIDAndOrderIDOutput{
		PickupID: pickupID,
		OrderID:  orderID,
	}
}
