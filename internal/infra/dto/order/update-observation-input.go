package orderdto

import orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"

type UpdateObservationOrder struct {
	Observation string `json:"observation"`
}

func (u *UpdateObservationOrder) UpdateModel(order *orderentity.Order) {
	order.Observation = u.Observation
}
