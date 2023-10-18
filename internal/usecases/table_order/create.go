package tableorderusecases

import (
	"context"

	"github.com/google/uuid"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) CreateTableOrder(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.CreateTableOrderInput) (uuid.UUID, error) {
	table, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.rt.CreateTableOrder(ctx, table); err != nil {
		return uuid.Nil, err
	}

	return table.ID, nil
}
