package orderusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) DeleteOrderByID(ctx context.Context, dtoId *entitydto.IdRequest) error {
	order, err := s.ro.GetOrderById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	if order.Status != orderentity.OrderStatusStaging {
		return errors.New("order must be staging")
	}

	return s.ro.DeleteOrder(ctx, dtoId.ID.String())
}
