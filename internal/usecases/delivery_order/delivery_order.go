package deliveryorderusecases

import (
	"context"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	deliveryorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

type IService interface {
	ICreateService
	IGetService
	IUpdateService
	IStatusService
}

type ICreateService interface {
	CreateDeliveryOrder(ctx context.Context, dto *deliveryorderdto.CreateDeliveryOrderInput) (*deliveryorderdto.DeliveryIDAndOrderIDOutput, error)
}

type IGetService interface {
	GetDeliveryById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.DeliveryOrder, error)
	GetAllDeliveries(ctx context.Context) ([]orderentity.DeliveryOrder, error)
	GetDeliveryOrderByStatus(ctx context.Context) (deliveries []orderentity.DeliveryOrder, err error)
	GetDeliveryOrderByClientId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.DeliveryOrder, error)
	GetDeliveryOrderByDriverId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.DeliveryOrder, error)
}

type IUpdateService interface {
	PendDeliveryOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	ShipDeliveryOrder(ctx context.Context, dtoID *entitydto.IdRequest, dtoDriver *deliveryorderdto.UpdateDriverOrder) (err error)
	DeliveryOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	UpdateDeliveryDriver(ctx context.Context, dto *entitydto.IdRequest, deliveryOrder *deliveryorderdto.UpdateDriverOrder) (err error)
}

type IStatusService interface {
	GetAllDeliveryOrderStatus(ctx context.Context) (deliveries []orderentity.StatusDeliveryOrder)
}

type Service struct {
	rdo orderentity.DeliveryOrderRepository
	ra  addressentity.Repository
	rc  cliententity.Repository
	ro  orderentity.OrderRepository
	re  employeeentity.Repository
	os  *orderusecases.Service
}

func NewService(rdo orderentity.DeliveryOrderRepository, ra addressentity.Repository, rc cliententity.Repository, ro orderentity.OrderRepository, re employeeentity.Repository, os *orderusecases.Service) IService {
	return &Service{rdo: rdo, ra: ra, rc: rc, ro: ro, re: re, os: os}
}
