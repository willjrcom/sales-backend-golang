package orderpickupusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
)

func (s *Service) GetPickupById(ctx context.Context, dto *entitydto.IDRequest) (*orderpickupdto.OrderPickupDTO, error) {
	if orderPickupModel, err := s.rp.GetPickupById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		pickup := orderPickupModel.ToDomain()
		orderPickupDTO := &orderpickupdto.OrderPickupDTO{}
		orderPickupDTO.FromDomain(pickup)
		return orderPickupDTO, nil
	}
}

func (s *Service) GetAllPickups(ctx context.Context) ([]orderpickupdto.OrderPickupDTO, error) {
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

func (s *Service) GetOrderPickupByStatus(ctx context.Context) (pickups []orderpickupdto.OrderPickupDTO, err error) {
	return nil, nil
}
