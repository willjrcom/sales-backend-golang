package orderdeliveryusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
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
	AddDependencies(ra model.AddressRepository, rc model.ClientRepository, ro model.OrderRepository, os *orderusecases.Service, rdd model.DeliveryDriverRepository)
}

type ICreateService interface {
	CreateOrderDelivery(ctx context.Context, dto *orderdeliverydto.DeliveryOrderCreateDTO) (*orderdeliverydto.OrderDeliveryIDDTO, error)
}

type IGetService interface {
	GetDeliveryById(ctx context.Context, dto *entitydto.IDRequest) (*orderentity.OrderDelivery, error)
	GetAllDeliveries(ctx context.Context) ([]orderentity.OrderDelivery, error)
	GetOrderDeliveryByStatus(ctx context.Context) (deliveries []orderentity.OrderDelivery, err error)
	GetOrderDeliveryByClientId(ctx context.Context, dto *entitydto.IDRequest) ([]orderentity.OrderDelivery, error)
	GetOrderDeliveryByDriverId(ctx context.Context, dto *entitydto.IDRequest) ([]orderentity.OrderDelivery, error)
}

type IUpdateService interface {
	PendOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	ShipOrderDelivery(ctx context.Context, dtoDriver *orderdeliverydto.DeliveryOrderUpdateShipDTO) (err error)
	OrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	UpdateDeliveryDriver(ctx context.Context, dto *entitydto.IDRequest, orderDelivery *orderdeliverydto.DeliveryOrderDriverUpdateDTO) (err error)
}

type IStatusService interface {
	GetAllOrderDeliveryStatus(ctx context.Context) (deliveries []orderentity.StatusOrderDelivery)
}

type Service struct {
	rdo model.OrderDeliveryRepository
	ra  model.AddressRepository
	rc  model.ClientRepository
	ro  model.OrderRepository
	os  *orderusecases.Service
	rdd model.DeliveryDriverRepository
}

func NewService(rdo model.OrderDeliveryRepository) IService {
	return &Service{rdo: rdo}
}

func (s *Service) AddDependencies(ra model.AddressRepository, rc model.ClientRepository, ro model.OrderRepository, os *orderusecases.Service, rdd model.DeliveryDriverRepository) {
	s.ra = ra
	s.rc = rc
	s.ro = ro
	s.os = os
	s.rdd = rdd
}
