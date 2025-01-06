package orderdeliverydto

import "github.com/google/uuid"

type OrderDeliveryIDDTO struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
	OrderID    uuid.UUID `json:"order_id"`
}

func FromDomain(deliveryID uuid.UUID, orderID uuid.UUID) *OrderDeliveryIDDTO {
	return &OrderDeliveryIDDTO{
		DeliveryID: deliveryID,
		OrderID:    orderID,
	}
}
