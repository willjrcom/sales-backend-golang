package orderdeliveryusecases

import (
	"context"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

type IService interface {
	ISetupService
	ICreateService
	IGetService
	IUpdateService
	IStatusService
}
type ISetupService interface {
	AddDependencies(ra addressentity.Repository, rc cliententity.Repository, ro orderentity.OrderRepository, re employeeentity.Repository, os *orderusecases.Service)
}

type ICreateService interface {
	CreateOrderDelivery(ctx context.Context, dto *orderdeliverydto.CreateOrderDeliveryInput) (*orderdeliverydto.DeliveryIDAndOrderIDOutput, error)
}

type IGetService interface {
	GetDeliveryById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.OrderDelivery, error)
	GetAllDeliveries(ctx context.Context) ([]orderentity.OrderDelivery, error)
	GetOrderDeliveryByStatus(ctx context.Context) (deliveries []orderentity.OrderDelivery, err error)
	GetOrderDeliveryByClientId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.OrderDelivery, error)
	GetOrderDeliveryByDriverId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.OrderDelivery, error)
}

type IUpdateService interface {
	PendOrderDelivery(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	ShipOrderDelivery(ctx context.Context, dtoID *entitydto.IdRequest, dtoDriver *orderdeliverydto.UpdateDriverOrder) (err error)
	OrderDelivery(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	UpdateDeliveryDriver(ctx context.Context, dto *entitydto.IdRequest, orderDelivery *orderdeliverydto.UpdateDriverOrder) (err error)
}

type IStatusService interface {
	GetAllOrderDeliveryStatus(ctx context.Context) (deliveries []orderentity.StatusOrderDelivery)
}

type Service struct {
	rdo orderentity.OrderDeliveryRepository
	ra  addressentity.Repository
	rc  cliententity.Repository
	ro  orderentity.OrderRepository
	re  employeeentity.Repository
	os  *orderusecases.Service
}

func NewService(rdo orderentity.OrderDeliveryRepository) IService {
	return &Service{rdo: rdo}
}

func (s *Service) AddDependencies(ra addressentity.Repository, rc cliententity.Repository, ro orderentity.OrderRepository, re employeeentity.Repository, os *orderusecases.Service) {
	s.ra = ra
	s.rc = rc
	s.ro = ro
	s.re = re
	s.os = os
}
