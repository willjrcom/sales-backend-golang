package pickuporderusecases

import (
	"context"

	"github.com/google/uuid"
	pickuporderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/pickup_order"
)

func (s *Service) CreatePickupOrder(ctx context.Context, dto *pickuporderdto.CreatePickupOrderInput) (uuid.UUID, error) {
	pickupOrder, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return uuid.Nil, err
	}

	pickupOrder.OrderID = orderID

	if err = s.rp.CreatePickupOrder(ctx, pickupOrder); err != nil {
		return uuid.Nil, err
	}

	return pickupOrder.ID, nil
}
