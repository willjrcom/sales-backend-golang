package ordertableusecases

import (
	"context"
	"errors"

	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
)

var (
	ErrTableIsNotAvailable = errors.New("table is not available")
)

func (s *Service) CreateOrderTable(ctx context.Context, dto *ordertabledto.CreateOrderTableInput) (*ordertabledto.TableIDAndOrderIDOutput, error) {
	orderTable, err := dto.ToModel()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	orderTable.OrderID = orderID

	table, err := s.rt.GetTableById(ctx, orderTable.TableID.String())

	if err != nil {
		return nil, err
	}

	table.LockTable()

	if orderTable.Name != "" {
		orderTable.Name = table.Name + " - " + orderTable.Name
	}

	if err = s.rto.CreateOrderTable(ctx, orderTable); err != nil {
		return nil, err
	}

	if err = s.rt.UpdateTable(ctx, table); err != nil {
		return nil, err
	}

	return ordertabledto.NewOutput(orderTable.ID, orderID), nil
}
