package deliveryorderusecases

import (
	"context"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

type IService interface {
	ICreateService
	IGetService
	IUpdateService
	IStatusService
}

type ICreateService interface {
	CreateDeliveryOrder(ctx context.Context, dto *orderdto.CreateDeliveryOrderInput) (uuid.UUID, error)
}

type IGetService interface {
	GetDeliveryById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.DeliveryOrder, error)
	GetAllDeliveries(ctx context.Context) ([]orderentity.DeliveryOrder, error)
	GetDeliveryOrderByStatus(ctx context.Context) (deliveries []orderentity.DeliveryOrder, err error)
	GetDeliveryOrderByClientId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.DeliveryOrder, error)
	GetDeliveryOrderByDriverId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.DeliveryOrder, error)
}

type IUpdateService interface {
	ShipDeliveryOrder(ctx context.Context, dtoId *entitydto.IdRequest, driverId *entitydto.IdRequest) (err error)
	DeliverDeliveryOrder(ctx context.Context, dtoId *entitydto.IdRequest) (err error)
	UpdateDeliveryOrder(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.DeliveryOrder) (err error)
	UpdateDeliveryAddress(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.UpdateDeliveryOrder) (err error)
	UpdateDriver(ctx context.Context, dto *entitydto.IdRequest, deliveryOrder *orderdto.UpdateDriverOrder) (err error)
}

type IStatusService interface {
	GetAllDeliveryOrderStatus(ctx context.Context) (deliveries []orderentity.StatusDeliveryOrder)
}

type Service struct {
	rdo orderentity.DeliveryRepository
	ra  addressentity.Repository
	rc  cliententity.Repository
	ro  orderentity.OrderRepository
	re  employeeentity.Repository
}

func NewService(rdo orderentity.DeliveryRepository, ra addressentity.Repository, rc cliententity.Repository, ro orderentity.OrderRepository, re employeeentity.Repository) IService {
	return &Service{rdo: rdo, ra: ra, rc: rc, ro: ro, re: re}
}
