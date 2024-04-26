package pickuporderusecases

import (
	"context"

	pickuporderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/pickup_order"
)

func (s *Service) CreatePickupOrder(ctx context.Context, dto *pickuporderdto.CreatePickupOrderInput) (*pickuporderdto.PickupIDAndOrderIDOutput, error) {
	pickupOrder, err := dto.ToModel()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	pickupOrder.OrderID = orderID

	if err = s.rp.CreatePickupOrder(ctx, pickupOrder); err != nil {
		return nil, err
	}

	return pickuporderdto.NewOutput(pickupOrder.ID, orderID), nil
}
