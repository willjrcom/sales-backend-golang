package orderpickupusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
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
	AddDependencies(os *orderusecases.Service)
}

type ICreateService interface {
	CreateOrderPickup(ctx context.Context, dto *orderpickupdto.OrderPickupCreateDTO) (*orderpickupdto.PickupIDAndOrderIDDTO, error)
}

type IGetService interface {
	GetPickupById(ctx context.Context, dto *entitydto.IDRequest) (*orderentity.OrderPickup, error)
	GetAllPickups(ctx context.Context) ([]orderentity.OrderPickup, error)
	GetOrderPickupByStatus(ctx context.Context) (pickups []orderentity.OrderPickup, err error)
}

type IUpdateService interface {
	PendingOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	ReadyOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	UpdateName(ctx context.Context, dtoID *entitydto.IDRequest, dtoPickup *orderpickupdto.UpdateOrderPickupInput) (err error)
}

type IStatusService interface {
	GetAllOrderPickupStatus(ctx context.Context) (pickups []orderentity.StatusOrderPickup)
}

type Service struct {
	rp model.OrderPickupRepository
	os *orderusecases.Service
}

func NewService(rp model.OrderPickupRepository) IService {
	return &Service{rp: rp}
}

func (s *Service) AddDependencies(os *orderusecases.Service) {
	s.os = os
}
