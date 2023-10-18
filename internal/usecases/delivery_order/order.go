package deliveryorderusecases

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Service struct {
	ro orderentity.DeliveryRepository
	ra addressentity.Repository
}

func NewService(ro orderentity.DeliveryRepository, ra addressentity.Repository) *Service {
	return &Service{ro: ro, ra: ra}
}
