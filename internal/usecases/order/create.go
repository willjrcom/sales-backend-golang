package orderusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *Service) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	shiftModel, err := s.rs.GetOpenedShift(ctx)
	if err != nil {
		return uuid.Nil, errors.New("must be opened shift")
	}

	shift := shiftModel.ToDomain()

	shift.IncrementCurrentOrder()

	shiftModel.FromDomain(shift)
	if err = s.rs.UpdateShift(ctx, shiftModel); err != nil {
		return uuid.Nil, err
	}

	order := orderentity.NewDefaultOrder(&shift.ID, shift.CurrentOrderNumber, shift.AttendantID)

	orderModel := &model.Order{}
	orderModel.FromDomain(order)

	if err := s.ro.CreateOrder(ctx, orderModel); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
