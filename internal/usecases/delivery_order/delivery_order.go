package deliveryorderusecases

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Service struct {
	rdo orderentity.DeliveryRepository
	ra  addressentity.Repository
	rc  cliententity.Repository
	ro  orderentity.OrderRepository
}

func NewService(rdo orderentity.DeliveryRepository, ra addressentity.Repository, rc cliententity.Repository, ro orderentity.OrderRepository) *Service {
	return &Service{rdo: rdo, ra: ra, rc: rc, ro: ro}
}
