package orderusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *OrderService) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	shiftModel, err := s.rs.GetCurrentShift(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("must open a new shift")
	}

	currentOrderNumber, err := s.rs.IncrementCurrentOrder(shiftModel.ID.String())
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
