package tableorderusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetTableById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.TableOrder, error) {
	if order, err := s.rto.GetTableOrderById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return order, nil
	}
}

func (s *Service) GetAllTables(ctx context.Context) ([]orderentity.TableOrder, error) {
	if orders, err := s.rto.GetAllTableOrders(ctx); err != nil {
		return nil, err
	} else {
		return orders, nil
	}
}
