package deliveryorderusecases

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Service struct {
	rdo orderentity.DeliveryRepository
	ra  addressentity.Repository
	rc  cliententity.Repository
	ro  orderentity.OrderRepository
	re  employeeentity.Repository
}

func NewService(rdo orderentity.DeliveryRepository, ra addressentity.Repository, rc cliententity.Repository, ro orderentity.OrderRepository, re employeeentity.Repository) *Service {
	return &Service{rdo: rdo, ra: ra, rc: rc, ro: ro, re: re}
}
