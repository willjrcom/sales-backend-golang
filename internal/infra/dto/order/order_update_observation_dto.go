package orderdto

import orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"

type OrderUpdateObservationDTO struct {
	Observation string `json:"observation"`
}

func (u *OrderUpdateObservationDTO) UpdateDomain(order *orderentity.Order) {
	order.Observation = u.Observation
}
