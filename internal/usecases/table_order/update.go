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

	if !newTable.IsAvailable && !dtoNew.ForceUpdate {
		return ErrTableNotAvailableToChange
	}

	tableOrder, err := s.rto.GetTableOrderById(ctx, dtoTableOrder.ID.String())

	if err != nil {
		return err
	}

	if tableOrder.TableID == newTable.ID {
		return errors.New("table order is already in this table")
	}

	table, err := s.rt.GetTableById(ctx, tableOrder.TableID.String())

	if err != nil {
		return err
	}

	tablesOrdersTogether, err := s.rto.GetPendingTableOrdersByTableId(ctx, tableOrder.TableID.String())
	if err != nil {
		return err
	}

	if len(tablesOrdersTogether) == 1 {
		table.UnlockTable()

		if err = s.rt.UpdateTable(ctx, table); err != nil {
			return err
		}
	}

	newTable.LockTable()

	if err = s.rt.UpdateTable(ctx, newTable); err != nil {
		return err
	}

	tableOrder.TableID = newTable.ID

	return s.rto.UpdateTableOrder(ctx, tableOrder)

}

func (s *Service) CloseTableOrder(ctx context.Context, dtoID *entitydto.IdRequest) error {
	tableOrder, err := s.rto.GetTableOrderById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := tableOrder.Close(); err != nil {
		return err
	}

	table, err := s.rt.GetTableById(ctx, tableOrder.TableID.String())
	if err != nil {
		return err
	}

	tablesOrdersTogether, err := s.rto.GetPendingTableOrdersByTableId(ctx, tableOrder.TableID.String())

	if err != nil {
		return err
	}

	if len(tablesOrdersTogether) == 1 {
		table.UnlockTable()
		if err := s.rt.UpdateTable(ctx, table); err != nil {
			return err
		}
	}

	return s.rto.UpdateTableOrder(ctx, tableOrder)
}
