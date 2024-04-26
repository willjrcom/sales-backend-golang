package tableorderusecases

import (
	"context"
	"errors"

	tableorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table_order"
)

var (
	ErrTableIsNotAvailable = errors.New("table is not available")
)

func (s *Service) CreateTableOrder(ctx context.Context, dto *tableorderdto.CreateTableOrderInput) (*tableorderdto.TableIDAndOrderIDOutput, error) {
	tableOrder, err := dto.ToModel()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	tableOrder.OrderID = orderID

	table, err := s.rt.GetTableById(ctx, tableOrder.TableID.String())

	if err != nil {
		return nil, err
	}

	table.LockTable()

	if tableOrder.Name != "" {
		tableOrder.Name = table.Name + " - " + tableOrder.Name
	}

	if err = s.rto.CreateTableOrder(ctx, tableOrder); err != nil {
		return nil, err
	}

	if err = s.rt.UpdateTable(ctx, table); err != nil {
		return nil, err
	}

	return tableorderdto.NewOutput(tableOrder.TableID, orderID), nil
}
