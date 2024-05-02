package orderusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

func (s *Service) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	shift, err := s.rs.GetOpenedShift(ctx)
	if err != nil {
		return uuid.Nil, errors.New("must be opened shift")
	}

	shift.IncrementCurrentOrder()
	if err = s.rs.UpdateShift(ctx, shift); err != nil {
		return uuid.Nil, err
	}

	order := orderentity.NewDefaultOrder(&shift.ID, shift.CurrentOrderNumber, shift.AttendantID)

	if err := s.ro.CreateOrder(ctx, order); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
