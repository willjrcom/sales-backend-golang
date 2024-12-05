package orderusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetOrderById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.Order, error) {
	if order, err := s.ro.GetOrderById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return order, nil
	}
}

func (s *Service) GetAllOrders(ctx context.Context) ([]orderentity.Order, error) {
	if orders, err := s.ro.GetAllOrders(ctx); err != nil {
		return nil, err
	} else {
		return orders, nil
	}
}

func (s *Service) GetAllOrderDeliveryStatus(ctx context.Context) ([]orderentity.Order, error) {
	return []orderentity.Order{}, nil
}
