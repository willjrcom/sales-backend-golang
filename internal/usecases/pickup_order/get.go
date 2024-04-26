package pickuporderusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetPickupById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.PickupOrder, error) {
	if pickupOrder, err := s.rp.GetPickupById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return pickupOrder, nil
	}
}

func (s *Service) GetAllPickups(ctx context.Context) ([]orderentity.PickupOrder, error) {
	if pickups, err := s.rp.GetAllPickups(ctx); err != nil {
		return nil, err
	} else {
		return pickups, nil
	}
}

func (s *Service) GetAllPickupOrderStatus(ctx context.Context) (pickups []orderentity.StatusPickupOrder) {
	return orderentity.GetAllPickupStatus()
}

func (s *Service) GetPickupOrderByStatus(ctx context.Context) (pickups []orderentity.PickupOrder, err error) {
	return nil, nil
}
