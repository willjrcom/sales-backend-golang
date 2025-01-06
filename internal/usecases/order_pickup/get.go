package orderpickupusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetPickupById(ctx context.Context, dto *entitydto.IDRequest) (*orderentity.OrderPickup, error) {
	if orderPickupModel, err := s.rp.GetPickupById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return orderPickupModel.ToDomain(), nil
	}
}

func (s *Service) GetAllPickups(ctx context.Context) ([]orderentity.OrderPickup, error) {
	if pickupModels, err := s.rp.GetAllPickups(ctx); err != nil {
		return nil, err
	} else {
		pickups := make([]orderentity.OrderPickup, len(pickupModels))
		for i, pickupModel := range pickupModels {
			pickups[i] = *pickupModel.ToDomain()
		}
		return pickups, nil
	}
}

func (s *Service) GetAllOrderPickupStatus(ctx context.Context) (pickups []orderentity.StatusOrderPickup) {
	return orderentity.GetAllPickupStatus()
}

func (s *Service) GetOrderPickupByStatus(ctx context.Context) (pickups []orderentity.OrderPickup, err error) {
	return nil, nil
}
