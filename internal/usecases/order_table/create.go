package ordertableusecases

import (
	"context"
	"errors"

	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrTableIsNotAvailable = errors.New("table is not available")
)

func (s *Service) CreateOrderTable(ctx context.Context, dto *ordertabledto.CreateOrderTableInput) (*ordertabledto.OrderTableIDDTO, error) {
	orderTable, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	orderTable.OrderID = orderID

	tableModel, err := s.rt.GetTableById(ctx, orderTable.TableID.String())

	if err != nil {
		return nil, err
	}

	table := tableModel.ToDomain()
	table.LockTable()

	if orderTable.Name == "" {
		orderTable.Name = table.Name
	} else {
		orderTable.Name = table.Name + " - " + orderTable.Name
	}

	company, err := s.cs.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	orderTable.UpdatePreferences(company.Preferences)

	orderTableModel := &model.OrderTable{}
	orderTableModel.FromDomain(orderTable)
	if err = s.rto.CreateOrderTable(ctx, orderTableModel); err != nil {
		return nil, err
	}

	tableModel.FromDomain(table)
	if err = s.rt.UpdateTable(ctx, tableModel); err != nil {
		return nil, err
	}

	return ordertabledto.FromDomain(orderTable.ID, orderID), nil
}
