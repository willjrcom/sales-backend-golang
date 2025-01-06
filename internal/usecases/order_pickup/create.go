package orderpickupusecases

import (
	"context"

	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
)

func (s *Service) CreateOrderPickup(ctx context.Context, dto *orderpickupdto.OrderPickupCreateDTO) (*orderpickupdto.PickupIDAndOrderIDDTO, error) {
	orderPickup, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	orderPickup.OrderID = orderID

	if err = s.rp.CreateOrderPickup(ctx, orderPickup); err != nil {
		return nil, err
	}

	return orderpickupdto.FromDomain(orderPickup.ID, orderID), nil
}
