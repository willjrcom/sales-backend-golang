package orderusecases

import (
	"context"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

func (s *Service) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	order := orderentity.NewDefaultOrder()

	// Get order Number
	order.OrderNumber = 1

	if err := s.ro.CreateOrder(ctx, order); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
