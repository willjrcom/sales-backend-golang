package orderrepositorylocal

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderRepositoryLocal(db *bun.DB) *OrderRepositoryBun {
	return &OrderRepositoryBun{db: db}
}

func (r *OrderRepositoryBun) CreateOrder(ctx context.Context, o *orderentity.Order) error {
	return nil
}

func (r *OrderRepositoryBun) UpdateOrder(ctx context.Context, o *orderentity.Order) error {
	return nil
}

func (r *OrderRepositoryBun) DeleteOrder(ctx context.Context, id string) error {
	return nil
}

func (r *OrderRepositoryBun) GetOrderById(ctx context.Context, id string) (*orderentity.Order, error) {
	return nil, nil
}

func (r *OrderRepositoryBun) GetOrderBy(ctx context.Context, o *orderentity.Order) (*orderentity.Order, error) {
	return nil, nil
}

func (r *OrderRepositoryBun) GetAllOrders(ctx context.Context) ([]orderentity.Order, error) {
	return nil, nil
}
