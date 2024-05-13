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

func (s *Service) ChangeTable(ctx context.Context, dtoOrderTable *entitydto.IdRequest, dtoNew *ordertabledto.UpdateOrderTableInput) error {
	newTable, err := s.rt.GetTableById(ctx, dtoNew.TableID.String())

	if err != nil {
		return err
	}

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

	table, err := s.rt.GetTableById(ctx, orderTable.TableID.String())

	if err != nil {
		return err
	}

	tablesOrdersTogether, err := s.rto.GetPendingOrderTablesByTableId(ctx, orderTable.TableID.String())
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

	orderTable.TableID = newTable.ID

	return s.rto.UpdateOrderTable(ctx, orderTable)

}

func (s *Service) CloseOrderTable(ctx context.Context, dtoID *entitydto.IdRequest) error {
	orderTable, err := s.rto.GetOrderTableById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := orderTable.Close(); err != nil {
		return err
	}

	table, err := s.rt.GetTableById(ctx, orderTable.TableID.String())
	if err != nil {
		return err
	}

	tablesOrdersTogether, err := s.rto.GetPendingOrderTablesByTableId(ctx, orderTable.TableID.String())

	if err != nil {
		return err
	}

	if len(tablesOrdersTogether) == 1 {
		table.UnlockTable()
		if err := s.rt.UpdateTable(ctx, table); err != nil {
			return err
		}
	}

	return s.rto.UpdateOrderTable(ctx, orderTable)
}
