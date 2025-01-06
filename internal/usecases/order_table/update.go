package ordertableusecases

import (
	"context"
	"errors"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
)

var (
	ErrTableNotAvailableToChange = errors.New("table not available to change")
)

func (s *Service) ChangeTable(ctx context.Context, dtoOrderTable *entitydto.IDRequest, dtoNew *ordertabledto.OrderTableUpdateDTO) error {
	newTableModel, err := s.rt.GetTableById(ctx, dtoNew.TableID.String())

	if err != nil {
		return err
	}

	newTable := newTableModel.ToDomain()

	if !newTable.IsAvailable && !dtoNew.ForceUpdate {
		return ErrTableNotAvailableToChange
	}

	orderTable, err := s.rto.GetOrderTableById(ctx, dtoOrderTable.ID.String())

	if err != nil {
		return err
	}

	if orderTable.TableID == newTable.ID {
		return errors.New("table order is already in this table")
	}

	tableModel, err := s.rt.GetTableById(ctx, orderTable.TableID.String())

	if err != nil {
		return err
	}

	table := tableModel.ToDomain()

	tablesOrdersTogether, err := s.rto.GetPendingOrderTablesByTableId(ctx, orderTable.TableID.String())
	if err != nil {
		return err
	}

	if len(tablesOrdersTogether) == 1 {
		table.UnlockTable()

		tableModel.FromDomain(table)
		if err = s.rt.UpdateTable(ctx, tableModel); err != nil {
			return err
		}
	}

	newTable.LockTable()

	newTableModel.FromDomain(newTable)
	if err = s.rt.UpdateTable(ctx, newTableModel); err != nil {
		return err
	}

	orderTable.TableID = newTable.ID

	return s.rto.UpdateOrderTable(ctx, orderTable)

}

func (s *Service) CloseOrderTable(ctx context.Context, dtoID *entitydto.IDRequest) error {
	orderTableModel, err := s.rto.GetOrderTableById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderTable := orderTableModel.ToDomain()

	if err := orderTable.Close(); err != nil {
		return err
	}

	tableModel, err := s.rt.GetTableById(ctx, orderTable.TableID.String())
	if err != nil {
		return err
	}

	table := tableModel.ToDomain()

	tablesOrdersTogether, err := s.rto.GetPendingOrderTablesByTableId(ctx, orderTable.TableID.String())

	if err != nil {
		return err
	}

	if len(tablesOrdersTogether) == 1 {
		table.UnlockTable()

		tableModel.FromDomain(table)
		if err := s.rt.UpdateTable(ctx, tableModel); err != nil {
			return err
		}
	}

	orderTableModel.FromDomain(orderTable)
	return s.rto.UpdateOrderTable(ctx, orderTableModel)
}
