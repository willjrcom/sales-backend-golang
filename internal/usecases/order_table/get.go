package ordertableusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetTableById(ctx context.Context, dto *entitydto.IDRequest) (*orderentity.OrderTable, error) {
	if order, err := s.rto.GetOrderTableById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return order, nil
	}
}

func (s *Service) GetAllTables(ctx context.Context) ([]orderentity.OrderTable, error) {
	if orders, err := s.rto.GetAllOrderTables(ctx); err != nil {
		return nil, err
	} else {
		return orders, nil
	}
}
