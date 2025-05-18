package ordertablerepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderTableRepositoryLocal struct {}

func NewOrderTableRepositoryLocal() model.OrderTableRepository {
	return &OrderTableRepositoryLocal{}
}

func (r *OrderTableRepositoryLocal) CreateOrderTable(ctx context.Context, table *model.OrderTable) error {
	return nil
}

func (r *OrderTableRepositoryLocal) UpdateOrderTable(ctx context.Context, table *model.OrderTable) error {
	return nil
}

func (r *OrderTableRepositoryLocal) DeleteOrderTable(ctx context.Context, id string) error {
	return nil
}

func (r *OrderTableRepositoryLocal) GetOrderTableById(ctx context.Context, id string) (*model.OrderTable, error) {
	return nil, nil
}

func (r *OrderTableRepositoryLocal) GetPendingOrderTablesByTableId(ctx context.Context, id string) ([]model.OrderTable, error) {
	return nil, nil
}

func (r *OrderTableRepositoryLocal) GetAllOrderTables(ctx context.Context) ([]model.OrderTable, error) {
	return nil, nil
}
