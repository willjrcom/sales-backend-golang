package orderusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetOrderById(ctx context.Context, dto *entitydto.IDRequest) (*orderentity.Order, error) {
	if orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {

		return orderModel.ToDomain(), nil
	}
}

func (s *Service) GetAllOrders(ctx context.Context) ([]orderentity.Order, error) {
	if orderModels, err := s.ro.GetAllOrders(ctx); err != nil {
		return nil, err
	} else {
		orders := make([]orderentity.Order, 0)
		for _, orderModel := range orderModels {
			orders = append(orders, *orderModel.ToDomain())
		}
		return orders, nil
	}
}

func (s *Service) GetAllOrdersWithDelivery(ctx context.Context) ([]orderentity.Order, error) {
	if orderModels, err := s.ro.GetAllOrdersWithDelivery(ctx); err != nil {
		return nil, err
	} else {
		orders := make([]orderentity.Order, 0)
		for _, orderModel := range orderModels {
			orders = append(orders, *orderModel.ToDomain())
		}
		return orders, nil
	}
}
