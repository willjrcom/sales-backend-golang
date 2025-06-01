package orderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) GetOrderById(ctx context.Context, dto *entitydto.IDRequest) (*orderdto.OrderDTO, error) {
	if orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		order := orderModel.ToDomain()

		orderDTO := &orderdto.OrderDTO{}
		orderDTO.FromDomain(order)
		return orderDTO, nil
	}
}

func (s *Service) GetAllOrders(ctx context.Context) ([]orderdto.OrderDTO, error) {
	if orderModels, err := s.ro.GetAllOrders(ctx); err != nil {
		return nil, err
	} else {
		orders := make([]orderdto.OrderDTO, 0)
		for _, orderModel := range orderModels {
			order := orderModel.ToDomain()
			orderDTO := &orderdto.OrderDTO{}
			orderDTO.FromDomain(order)
			orders = append(orders, *orderDTO)
		}
		return orders, nil
	}
}

func (s *Service) GetAllOrdersWithDelivery(ctx context.Context, page, perPage int) ([]orderdto.OrderDTO, error) {
	if orderModels, err := s.ro.GetAllOrdersWithDelivery(ctx, page, perPage); err != nil {
		return nil, err
	} else {
		orders := make([]orderdto.OrderDTO, 0)
		for _, orderModel := range orderModels {
			order := orderModel.ToDomain()
			orderDTO := &orderdto.OrderDTO{}
			orderDTO.FromDomain(order)
			orders = append(orders, *orderDTO)
		}
		return orders, nil
	}
}
