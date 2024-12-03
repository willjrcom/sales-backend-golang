package orderpickupusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
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
	AddDependencies(os *orderusecases.Service)
}

type ICreateService interface {
	CreateOrderPickup(ctx context.Context, dto *orderpickupdto.CreateOrderPickupInput) (*orderpickupdto.PickupIDAndOrderIDOutput, error)
}

type IGetService interface {
	GetPickupById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.OrderPickup, error)
	GetAllPickups(ctx context.Context) ([]orderentity.OrderPickup, error)
	GetOrderPickupByStatus(ctx context.Context) (pickups []orderentity.OrderPickup, err error)
}

type IUpdateService interface {
	PendingOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	ReadyOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error)
	UpdateName(ctx context.Context, dtoID *entitydto.IdRequest, dtoPickup *orderpickupdto.UpdateOrderPickupInput) (err error)
}

type IStatusService interface {
	GetAllOrderPickupStatus(ctx context.Context) (pickups []orderentity.StatusOrderPickup)
}

type Service struct {
	rp orderentity.OrderPickupRepository
	os *orderusecases.Service
}

func NewService(rp orderentity.OrderPickupRepository) IService {
	return &Service{rp: rp}
}

func (s *Service) AddDependencies(os *orderusecases.Service) {
	s.os = os
}
