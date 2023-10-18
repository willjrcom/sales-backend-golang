package orderrepositorybun

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

func NewOrderRepositoryBun(db *bun.DB) *OrderRepositoryBun {
	return &OrderRepositoryBun{db: db}
}

func (r *OrderRepositoryBun) CreateOrder(ctx context.Context, order *orderentity.Order) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(order).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) UpdateOrder(ctx context.Context, order *orderentity.Order) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(order).Where("id = ?", order.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) DeleteOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&orderentity.Order{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) GetOrderById(ctx context.Context, id string) (*orderentity.Order, error) {
	order := &orderentity.Order{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(order).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()
	// .Relation("Delivery").Relation("Table").Relation("Groups").Relation("Attendant")

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepositoryBun) GetOrderBy(ctx context.Context, order *orderentity.Order) ([]orderentity.Order, error) {
	orders := []orderentity.Order{}

	r.mu.Lock()
	query := r.db.NewSelect().Model(&orderentity.Order{})

	if order.Status != "" {
		query.Where("order.status = ?", order.Status)
	}

	err := query.Relation("Delivery").Relation("Table").Relation("Groups").Relation("Attendant").Scan(ctx, &orders)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepositoryBun) GetAllOrders(ctx context.Context) ([]orderentity.Order, error) {
	orders := []orderentity.Order{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&orders).Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return orders, nil
}
