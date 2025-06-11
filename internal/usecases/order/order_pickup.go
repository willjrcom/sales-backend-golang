package orderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type IPickupService interface {
	ISetupPickupService
	ICreatePickupService
	IGetPickupService
	IUpdatePickupService
}

type ISetupPickupService interface {
	AddDependencies(os *OrderService)
}

type ICreatePickupService interface {
	CreateOrderPickup(ctx context.Context, dto *orderpickupdto.OrderPickupCreateDTO) (*orderpickupdto.PickupIDAndOrderIDDTO, error)
}

type IGetPickupService interface {
	GetPickupById(ctx context.Context, dto *entitydto.IDRequest) (*orderpickupdto.OrderPickupDTO, error)
	GetAllPickups(ctx context.Context) ([]orderpickupdto.OrderPickupDTO, error)
}

type IUpdatePickupService interface {
	PendingOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	ReadyOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	CancelOrderPickup(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	UpdateName(ctx context.Context, dtoID *entitydto.IDRequest, dtoPickup *orderpickupdto.UpdateOrderPickupInput) (err error)
}

type OrderPickupService struct {
	rp model.OrderPickupRepository
	os *OrderService
}

func NewPickupService(rp model.OrderPickupRepository) IPickupService {
	return &OrderPickupService{rp: rp}
}

func (s *OrderPickupService) AddDependencies(os *OrderService) {
	s.os = os
}

func (s *OrderPickupService) CreateOrderPickup(ctx context.Context, dto *orderpickupdto.OrderPickupCreateDTO) (*orderpickupdto.PickupIDAndOrderIDDTO, error) {
	orderPickup, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	orderPickup.OrderID = orderID

	orderPickupModel := &model.OrderPickup{}
	orderPickupModel.FromDomain(orderPickup)

	if err = s.rp.CreateOrderPickup(ctx, orderPickupModel); err != nil {
		return nil, err
	}

	return orderpickupdto.FromDomain(orderPickup.ID, orderID), nil
}

func (s *OrderPickupService) PendingOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
	orderPickupModel, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderPickup := orderPickupModel.ToDomain()
	if err := orderPickup.Pend(); err != nil {
		return err
	}

	orderPickupModel.FromDomain(orderPickup)
	if err = s.rp.UpdateOrderPickup(ctx, orderPickupModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderPickupService) ReadyOrder(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
	orderPickupModel, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderPickup := orderPickupModel.ToDomain()
	if err := orderPickup.Ready(); err != nil {
		return err
	}

	orderPickupModel.FromDomain(orderPickup)
	if err = s.rp.UpdateOrderPickup(ctx, orderPickupModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderPickupService) CancelOrderPickup(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
	orderPickupModel, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderPickup := orderPickupModel.ToDomain()
	if err := orderPickup.Cancel(); err != nil {
		return err
	}

	orderPickupModel.FromDomain(orderPickup)
	if err = s.rp.UpdateOrderPickup(ctx, orderPickupModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderPickupService) UpdateName(ctx context.Context, dtoID *entitydto.IDRequest, dtoPickup *orderpickupdto.UpdateOrderPickupInput) (err error) {
	orderPickupModel, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderPickup := orderPickupModel.ToDomain()
	if err := orderPickup.UpdateName(dtoPickup.Name); err != nil {
		return err
	}

	orderPickupModel.FromDomain(orderPickup)
	if err = s.rp.UpdateOrderPickup(ctx, orderPickupModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderPickupService) GetPickupById(ctx context.Context, dto *entitydto.IDRequest) (*orderpickupdto.OrderPickupDTO, error) {
	if orderPickupModel, err := s.rp.GetPickupById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		pickup := orderPickupModel.ToDomain()
		orderPickupDTO := &orderpickupdto.OrderPickupDTO{}
		orderPickupDTO.FromDomain(pickup)
		return orderPickupDTO, nil
	}
}

func (s *OrderPickupService) GetAllPickups(ctx context.Context) ([]orderpickupdto.OrderPickupDTO, error) {
	if pickupModels, err := s.rp.GetAllPickups(ctx); err != nil {
		return nil, err
	} else {
		pickupDTOs := make([]orderpickupdto.OrderPickupDTO, 0)
		for _, pickupModel := range pickupModels {
			pickup := pickupModel.ToDomain()
			pickupDTO := &orderpickupdto.OrderPickupDTO{}
			pickupDTO.FromDomain(pickup)
			pickupDTOs = append(pickupDTOs, *pickupDTO)
		}
		return pickupDTOs, nil
	}
}
