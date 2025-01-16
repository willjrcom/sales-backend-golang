package orderusecases

import (
	"context"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *Service) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	shift, err := s.rs.GetCurrentShift(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	currentOrderNumber, err := s.rs.IncrementCurrentOrder(shift.ID.String())
	if err != nil {
		return uuid.Nil, err
	}

	order := orderentity.NewDefaultOrder(shift.ID, currentOrderNumber, shift.AttendantID)

	orderModel := &model.Order{}
	orderModel.FromDomain(order)

	if err := s.ro.CreateOrder(ctx, orderModel); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
