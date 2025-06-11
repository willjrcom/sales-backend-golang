package orderpickupusecases

import (
	"context"

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
}

type ISetupService interface {
	AddDependencies(os *orderusecases.OrderService)
}

type ICreateService interface {
	CreateOrderPickup(ctx context.Context, dto *orderpickupdto.OrderPickupCreateDTO) (*orderpickupdto.PickupIDAndOrderIDDTO, error)
}

type IGetService interface {
	GetPickupById(ctx context.Context, dto *entitydto.IDRequest) (*orderpickupdto.OrderPickupDTO, error)
	GetAllPickups(ctx context.Context) ([]orderpickupdto.OrderPickupDTO, error)
	GetOrderPickupByStatus(ctx context.Context) (pickups []orderpickupdto.OrderPickupDTO, err error)
}

type IUpdateService interface {
	PendingOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	ReadyOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	UpdateName(ctx context.Context, dtoID *entitydto.IDRequest, dtoPickup *orderpickupdto.UpdateOrderPickupInput) (err error)
}

type Service struct {
	rp model.OrderPickupRepository
	os *orderusecases.OrderService
}

func NewService(rp model.OrderPickupRepository) IService {
	return &Service{rp: rp}
}

func (s *Service) AddDependencies(os *orderusecases.OrderService) {
	s.os = os
}
