package deliveryorderusecases

import (
	"context"

	"github.com/google/uuid"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) CreateDeliveryOrder(ctx context.Context, dto *orderdto.CreateDeliveryOrderInput) (uuid.UUID, error) {
	delivery, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.ro.CreateDeliveryOrder(ctx, delivery); err != nil {
		return uuid.Nil, err
	}

	return delivery.ID, nil
}
