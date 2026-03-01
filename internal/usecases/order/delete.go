package orderusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *OrderService) DeleteOrderByID(ctx context.Context, dtoId *entitydto.IDRequest) error {
	order, err := s.ro.GetOrderById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}
	if order.Status != string(orderentity.OrderStatusStaging) {
		return errors.New("order is not in staging status")
	}

	orderDomain := order.ToDomain()
	if err := s.restoreStockFromOrder(ctx, orderDomain); err != nil {
		return err
	}

	return s.ro.DeleteOrder(ctx, dtoId.ID.String())
}
