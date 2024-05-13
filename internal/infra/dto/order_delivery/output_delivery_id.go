package orderdeliverydto

import "github.com/google/uuid"

type DeliveryIDAndOrderIDOutput struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
	OrderID    uuid.UUID `json:"order_id"`
}

func NewOutput(deliveryID uuid.UUID, orderID uuid.UUID) *DeliveryIDAndOrderIDOutput {
	return &DeliveryIDAndOrderIDOutput{
		DeliveryID: deliveryID,
		OrderID:    orderID,
	}
}
