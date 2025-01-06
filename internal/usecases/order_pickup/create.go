package orderpickupusecases

import (
	"context"

	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
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

	orderPickupModel := &model.OrderPickup{}
	orderPickupModel.FromDomain(orderPickup)

	if err = s.rp.CreateOrderPickup(ctx, orderPickupModel); err != nil {
		return nil, err
	}

	return orderpickupdto.FromDomain(orderPickup.ID, orderID), nil
}
