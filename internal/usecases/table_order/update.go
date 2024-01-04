package tableorderusecases

import (
	"context"
	"errors"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	tableorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table_order"
)

var (
	ErrTableNotAvailableToChange = errors.New("table not available to change")
)

func (s *Service) ChangeTable(ctx context.Context, dtoTableOrder *entitydto.IdRequest, dtoNew *tableorderdto.UpdateTableOrderInput) error {
	newTable, err := s.rt.GetTableById(ctx, dtoNew.TableID.String())

	if err != nil {
		return err
	}

	if !newTable.IsAvailable {
		return ErrTableNotAvailableToChange
	}

	tableOrder, err := s.rto.GetTableOrderById(ctx, dtoTableOrder.ID.String())

	if err != nil {
		return err
	}

	table, err := s.rt.GetTableById(ctx, tableOrder.TableID.String())

	if err != nil {
		return err
	}

	table.UnlockTable()

	if err = s.rt.UpdateTable(ctx, table); err != nil {
		return err
	}

	newTable.LockTable()

	if err = s.rt.UpdateTable(ctx, newTable); err != nil {
		return err
	}

	tableOrder.TableID = dtoNew.TableID

	return s.rto.UpdateTableOrder(ctx, tableOrder)

}

func (s *Service) FinishTableOrder(ctx context.Context, dtoID *entitydto.IdRequest) error {
	table, err := s.rt.GetTableById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	table.UnlockTable()

	return s.rt.UpdateTable(ctx, table)
}
