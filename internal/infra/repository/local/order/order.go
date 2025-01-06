package orderrepositorylocal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderRepositoryLocal struct {
	orders map[uuid.UUID]*model.Order
}

func NewOrderRepositoryLocal() *OrderRepositoryLocal {
	return &OrderRepositoryLocal{orders: make(map[uuid.UUID]*model.Order)}
}

func (r *OrderRepositoryLocal) CreateOrder(ctx context.Context, order *model.Order) error {
	if _, ok := r.orders[order.Entity.ID]; ok {
		return errors.New("order already exists")
	}

	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepositoryLocal) DeleteOrder(ctx context.Context, id string) error {
	if _, ok := r.orders[uuid.MustParse(id)]; !ok {

		return errors.New("order not found")
	}

	delete(r.orders, uuid.MustParse(id))
	return nil
}

func (r *OrderRepositoryLocal) PendingOrder(ctx context.Context, order *model.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepositoryLocal) UpdateOrder(ctx context.Context, order *model.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepositoryLocal) GetOrderById(ctx context.Context, id string) (*model.Order, error) {
	if p, ok := r.orders[uuid.MustParse(id)]; ok {
		return p, nil
	}

	return nil, errors.New("order not found")
}

func (r *OrderRepositoryLocal) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	orders := make([]model.Order, 0)

	for _, p := range r.orders {
		orders = append(orders, *p)
	}

	return orders, nil
}

func (r *OrderRepositoryLocal) UpdateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {
	return nil
}

func (r *OrderRepositoryLocal) DeleteOrderDelivery(ctx context.Context, id string) error {
	return nil
}

func (r *OrderRepositoryLocal) UpdateOrderTable(ctx context.Context, table *model.OrderTable) error {
	return nil
}

func (r *OrderRepositoryLocal) DeleteOrderTable(ctx context.Context, id string) error {
	return nil
}

func (r *OrderRepositoryLocal) AddPaymentOrder(ctx context.Context, payment *model.PaymentOrder) error {
	return nil
}
