package orderusecases

import (
	"context"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
)

func (s *OrderService) GetClientByID(ctx context.Context, dto *entitydto.IDRequest) (*clientdto.ClientDTO, error) {
	return s.clientService.GetClientById(ctx, dto)
}

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

func (s *OrderService) GetOrderIDFromOrderDeliveriesByClientId(ctx context.Context, dto *entitydto.IDRequest) ([]orderdto.OrderDTO, error) {
	deliveryModels, err := s.rdo.GetOrderIDFromOrderDeliveriesByClientId(ctx, dto.ID.String())
	if err != nil {
		return nil, err
	} else {
		orders := []orderdto.OrderDTO{}
		for _, deliveryModel := range deliveryModels {
			orderModel, err := s.ro.GetOrderById(ctx, deliveryModel.OrderID.String())
			if err != nil {
				continue
			}

			order := orderModel.ToDomain()

			dto := &orderdto.OrderDTO{}
			dto.FromDomain(order)
			orders = append(orders, *dto)
		}
		return orders, nil
	}
}

func (s *OrderService) GetOrdersTableByTableId(ctx context.Context, dto *ordertabledto.OrderTableContactInput) ([]orderdto.OrderDTO, error) {
	if tableModels, err := s.st.rto.GetOrderTablesByTableId(ctx, dto.TableID.String(), dto.Contact); err != nil {
		return nil, err
	} else {
		orders := []orderdto.OrderDTO{}
		for _, tableModel := range tableModels {
			orderModel, err := s.ro.GetOrderById(ctx, tableModel.OrderID.String())
			if err != nil {
				continue
			}

			order := orderModel.ToDomain()

			dto := &orderdto.OrderDTO{}
			dto.FromDomain(order)
			orders = append(orders, *dto)
		}
		return orders, nil
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

func (s *OrderService) GetOrdersPickupByContact(ctx context.Context, contact string) ([]orderdto.OrderDTO, error) {
	if pickups, err := s.sp.GetOrderIDFromOrderPickupsByContact(ctx, contact); err != nil {
		return nil, err
	} else {
		orders := []orderdto.OrderDTO{}
		for _, pickup := range pickups {
			orderModel, err := s.ro.GetOrderById(ctx, pickup.OrderID.String())
			if err != nil {
				continue
			}

			order := orderModel.ToDomain()

			dto := &orderdto.OrderDTO{}
			dto.FromDomain(order)
			orders = append(orders, *dto)
		}
		return orders, nil
	}
}

func (s *OrderService) GetAllOrdersWithPickup(ctx context.Context, status orderentity.StatusOrderPickup, page, perPage int) ([]orderdto.OrderDTO, error) {
	shiftID := uuid.Nil
	shiftModel, _ := s.rs.GetCurrentShift(ctx)
	if shiftModel != nil {
		shiftID = shiftModel.ID
	}

	if orderModels, err := s.ro.GetAllOrdersWithPickup(ctx, shiftID.String(), status, page, perPage); err != nil {
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
