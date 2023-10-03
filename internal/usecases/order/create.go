package orderusecases

import (
	"context"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	order := orderentity.NewDefaultOrder()

	// Get order Number
	order.OrderNumber = 1

	if err := s.ro.CreateOrder(ctx, order); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}

func (s *Service) CreateDeliveryOrder(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.CreateDeliveryOrderInput) (uuid.UUID, error) {
	order, err := s.ro.GetOrderById(ctx, dtoId.ID.String())

	if err != nil {
		return uuid.Nil, err
	}

	delivery, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.ro.UpdateOrder(ctx, order); err != nil {
		return uuid.Nil, err
	}

	if err = s.ro.UpdateDeliveryOrder(ctx, order, delivery); err != nil {
		return uuid.Nil, err
	}

	return order.Delivery.ID, nil
}

func (s *Service) CreateTableOrder(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.CreateTableOrderInput) (uuid.UUID, error) {
	order, err := s.ro.GetOrderById(ctx, dtoId.ID.String())

	if err != nil {
		return uuid.Nil, err
	}

	table, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.ro.UpdateTableOrder(ctx, order, table); err != nil {
		return uuid.Nil, err
	}

	return order.Table.ID, nil
}
