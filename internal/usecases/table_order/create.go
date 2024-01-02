package tableorderusecases

import (
	"context"

	"github.com/google/uuid"
	tableorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table_order"
)

func (s *Service) CreateTableOrder(ctx context.Context, dto *tableorderdto.CreateTableOrderInput) (uuid.UUID, error) {
	table, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if _, err = s.rt.GetTableById(ctx, table.TableID.String()); err != nil {
		return uuid.Nil, err
	}

	if err = s.rto.CreateTableOrder(ctx, table); err != nil {
		return uuid.Nil, err
	}

	return table.ID, nil
}
