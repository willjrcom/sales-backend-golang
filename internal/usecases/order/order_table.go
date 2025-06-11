package orderusecases

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

type OrderTableService struct {
	rto model.OrderTableRepository
	rt  model.TableRepository
	os  *OrderService
	cs  *companyusecases.Service
}

func NewOrderTableService(rto model.OrderTableRepository) *OrderTableService {
	return &OrderTableService{rto: rto}
}

func (s *OrderTableService) AddDependencies(rt model.TableRepository, os *OrderService, cs *companyusecases.Service) {
	s.rt = rt
	s.os = os
	s.cs = cs
}

var (
	ErrTableNotAvailableToChange = errors.New("table not available to change")
	ErrTableIsNotAvailable       = errors.New("table is not available")
)

func (s *OrderTableService) CreateOrderTable(ctx context.Context, dto *ordertabledto.CreateOrderTableInput) (*ordertabledto.OrderTableIDDTO, error) {
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

func (s *OrderTableService) GetTableById(ctx context.Context, dto *entitydto.IDRequest) (*ordertabledto.OrderTableDTO, error) {
	if orderTableModel, err := s.rto.GetOrderTableById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		orderTable := orderTableModel.ToDomain()

		orderTableDTO := &ordertabledto.OrderTableDTO{}
		orderTableDTO.FromDomain(orderTable)
		return orderTableDTO, nil
	}
}

func (s *OrderTableService) GetAllTables(ctx context.Context) ([]ordertabledto.OrderTableDTO, error) {
	if orderTableModels, err := s.rto.GetAllOrderTables(ctx); err != nil {
		return nil, err
	} else {
		orderTableDTOs := make([]ordertabledto.OrderTableDTO, 0)
		for _, orderTableModel := range orderTableModels {
			orderTable := orderTableModel.ToDomain()

			orderTableModelDTO := &ordertabledto.OrderTableDTO{}
			orderTableModelDTO.FromDomain(orderTable)
			orderTableDTOs = append(orderTableDTOs, *orderTableModelDTO)
		}
		return orderTableDTOs, nil
	}
}

func (s *OrderTableService) ChangeTable(ctx context.Context, dtoOrderTable *entitydto.IDRequest, dtoNew *ordertabledto.OrderTableUpdateDTO) error {
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

// AddTableTax applies the configured table tax rate to the order-table.
func (s *OrderTableService) AddTableTax(ctx context.Context, dtoID *entitydto.IDRequest) error {
	// Retrieve existing order-table record
	orderTableModel, err := s.rto.GetOrderTableById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}
	// Convert to domain
	orderTable := orderTableModel.ToDomain()
	// Get company preferences
	companyDTO, err := s.cs.GetCompany(ctx)
	if err != nil {
		return err
	}
	// Update tax rate based on preferences
	orderTable.UpdatePreferences(companyDTO.Preferences)

	// Persist changes
	orderTableModel.FromDomain(orderTable)
	if err := s.rto.UpdateOrderTable(ctx, orderTableModel); err != nil {
		return err
	}

	return s.os.UpdateOrderTotal(ctx, orderTable.OrderID.String())
}

// RemoveTableTax sets the table tax rate to zero for the order-table.
func (s *OrderTableService) RemoveTableTax(ctx context.Context, dtoID *entitydto.IDRequest) error {
	// Retrieve existing order-table record
	orderTableModel, err := s.rto.GetOrderTableById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	// Convert to domain and clear tax rate
	orderTable := orderTableModel.ToDomain()
	orderTable.TaxRate = decimal.Zero

	// Persist changes
	orderTableModel.FromDomain(orderTable)
	if err := s.rto.UpdateOrderTable(ctx, orderTableModel); err != nil {
		return err
	}

	return s.os.UpdateOrderTotal(ctx, orderTable.OrderID.String())
}

func (s *OrderTableService) CloseOrderTable(ctx context.Context, dtoID *entitydto.IDRequest) error {
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

func (s *OrderTableService) CancelOrderTable(ctx context.Context, dtoID *entitydto.IDRequest) error {
	orderTableModel, err := s.rto.GetOrderTableById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderTable := orderTableModel.ToDomain()

	if err := orderTable.Cancel(); err != nil {
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
