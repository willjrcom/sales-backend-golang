package pickuporderusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	pickuporderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/pickup_order"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

type IService interface {
	ICreateService
	IGetService
	IUpdateService
	IStatusService
}

type ICreateService interface {
	CreatePickupOrder(ctx context.Context, dto *pickuporderdto.CreatePickupOrderInput) (*pickuporderdto.PickupIDAndOrderIDOutput, error)
}

type IGetService interface {
	GetPickupById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.PickupOrder, error)
	GetAllPickups(ctx context.Context) ([]orderentity.PickupOrder, error)
	GetPickupOrderByStatus(ctx context.Context) (pickups []orderentity.PickupOrder, err error)
}

type IUpdateService interface {
	PendingOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	ReadyOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
}

type IStatusService interface {
	GetAllPickupOrderStatus(ctx context.Context) (pickups []orderentity.StatusPickupOrder)
}

type Service struct {
	rp orderentity.PickupOrderRepository
	os *orderusecases.Service
}

func NewService(rp orderentity.PickupOrderRepository, os *orderusecases.Service) IService {
	return &Service{rp: rp, os: os}
}
