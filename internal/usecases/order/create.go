package orderusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *OrderService) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	if s.sc != nil {
		if err := s.sc.ValidateSubscription(ctx); err != nil {
			return uuid.Nil, err
		}
	}

	shiftModel, err := s.rs.GetCurrentShift(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("must open a new shift")
	}

	currentOrderNumber, err := s.rs.IncrementCurrentOrder(ctx, shiftModel.ID.String())
	if err != nil {
		return uuid.Nil, err
	}

	order := orderentity.NewDefaultOrder(shiftModel.ID, currentOrderNumber, shiftModel.AttendantID)

	orderModel := &model.Order{}
	orderModel.FromDomain(order)

	if err := s.ro.CreateOrder(ctx, orderModel); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
