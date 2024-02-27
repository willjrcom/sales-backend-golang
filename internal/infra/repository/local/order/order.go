package orderrepositorylocal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderRepositoryLocal struct {
	orders map[uuid.UUID]*orderentity.Order
}

func NewOrderRepositoryLocal() *OrderRepositoryLocal {
	return &OrderRepositoryLocal{orders: make(map[uuid.UUID]*orderentity.Order)}
}

func (r *OrderRepositoryLocal) CreateOrder(ctx context.Context, order *orderentity.Order) error {
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

func (r *OrderRepositoryLocal) PendingOrder(ctx context.Context, order *orderentity.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepositoryLocal) UpdateOrder(ctx context.Context, order *orderentity.Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepositoryLocal) GetOrderById(ctx context.Context, id string) (*orderentity.Order, error) {
	if p, ok := r.orders[uuid.MustParse(id)]; ok {
		return p, nil
	}

	return nil, errors.New("order not found")
}

func (r *OrderRepositoryLocal) GetAllOrders(ctx context.Context) ([]orderentity.Order, error) {
	orders := make([]orderentity.Order, 0)

	for _, p := range r.orders {
		orders = append(orders, *p)
	}

	return orders, nil
}

func (r *OrderRepositoryLocal) UpdateDeliveryOrder(ctx context.Context, delivery *orderentity.DeliveryOrder) error {
	return nil
}

func (r *OrderRepositoryLocal) DeleteDeliveryOrder(ctx context.Context, id string) error {
	return nil
}

func (r *OrderRepositoryLocal) UpdateTableOrder(ctx context.Context, table *orderentity.TableOrder) error {
	return nil
}

func (r *OrderRepositoryLocal) DeleteTableOrder(ctx context.Context, id string) error {
	return nil
}

func (r *OrderRepositoryLocal) AddPaymentOrder(ctx context.Context, payment *orderentity.PaymentOrder) error {
	return nil
}
