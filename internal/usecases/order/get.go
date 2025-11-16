package orderusecases

import (
	"context"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *OrderService) GetOrderById(ctx context.Context, dto *entitydto.IDRequest) (*orderdto.OrderDTO, error) {
	if orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		order := orderModel.ToDomain()

		orderDTO := &orderdto.OrderDTO{}
		orderDTO.FromDomain(order)
		return orderDTO, nil
	}
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]orderdto.OrderDTO, error) {
	shiftID := uuid.Nil
	shiftModel, _ := s.rs.GetCurrentShift(ctx)
	if shiftModel != nil {
		shiftID = shiftModel.ID
	}

	validStatuses := []orderentity.StatusOrder{
		orderentity.OrderStatusStaging,
		orderentity.OrderStatusPending,
		orderentity.OrderStatusReady,
	}

	if orderModels, err := s.ro.GetAllOrders(ctx, shiftID.String(), validStatuses, false, "OR"); err != nil {
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

func (s *OrderService) GetAllOrdersWithDelivery(ctx context.Context, page, perPage int) ([]orderdto.OrderDTO, error) {
	shiftID := uuid.Nil
	shiftModel, _ := s.rs.GetCurrentShift(ctx)
	if shiftModel != nil {
		shiftID = shiftModel.ID
	}

	if orderModels, err := s.ro.GetAllOrdersWithDelivery(ctx, shiftID.String(), page, perPage); err != nil {
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

func (s *OrderService) GetAllOrdersWithPickup(ctx context.Context, page, perPage int) ([]orderdto.OrderDTO, error) {
	shiftID := uuid.Nil
	shiftModel, _ := s.rs.GetCurrentShift(ctx)
	if shiftModel != nil {
		shiftID = shiftModel.ID
	}

	if orderModels, err := s.ro.GetAllOrdersWithPickup(ctx, shiftID.String(), page, perPage); err != nil {
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
