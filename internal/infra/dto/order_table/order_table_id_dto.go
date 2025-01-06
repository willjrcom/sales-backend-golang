package ordertabledto

import "github.com/google/uuid"

type OrderTableIDDTO struct {
	TableID uuid.UUID `json:"table_id"`
	OrderID uuid.UUID `json:"order_id"`
}

func FromDomain(tableID uuid.UUID, orderID uuid.UUID) *OrderTableIDDTO {
	return &OrderTableIDDTO{
		TableID: tableID,
		OrderID: orderID,
	}
}
