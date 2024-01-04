package tableorderusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	tableorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table_order"
)

var (
	ErrTableIsNotAvailable = errors.New("table is not available")
)

func (s *Service) CreateTableOrder(ctx context.Context, dto *tableorderdto.CreateTableOrderInput) (uuid.UUID, error) {
	tableOrder, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	table, err := s.rt.GetTableById(ctx, tableOrder.TableID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if !table.IsAvailable {
		return uuid.Nil, ErrTableIsNotAvailable
	}

	table.LockTable()

	tableOrder.Name = table.Name

	if err = s.rto.CreateTableOrder(ctx, tableOrder); err != nil {
		return uuid.Nil, err
	}

	if err = s.rt.UpdateTable(ctx, table); err != nil {
		return uuid.Nil, err
	}

	return tableOrder.ID, nil
}
