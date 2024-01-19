package orderusecases

import (
	"context"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) CreateDefaultOrder(ctx context.Context, dto *orderdto.CreateOrderInput) (uuid.UUID, error) {
	shiftID, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, err
	}

	shift, err := s.rs.GetShiftByID(ctx, shiftID.String())
	if err != nil {
		return uuid.Nil, err
	}

	shift.IncrementCurrentOrder()
	if err = s.rs.UpdateShift(ctx, shift); err != nil {
		return uuid.Nil, err
	}

	order := orderentity.NewDefaultOrder(shiftID, shift.CurrentOrderNumber, shift.AttendantID)

	if err := s.ro.CreateOrder(ctx, order); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
