package orderrepositorylocal

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderRepositoryLocal struct {
	mu     sync.Mutex
	orders map[uuid.UUID]*orderentity.Order
}

func NewOrderRepositoryLocal() *OrderRepositoryLocal {
	return &OrderRepositoryLocal{orders: make(map[uuid.UUID]*orderentity.Order)}
}

func (r *OrderRepositoryLocal) CreateOrder(order *orderentity.Order) error {
	r.mu.Lock()

	if _, ok := r.orders[order.Entity.ID]; ok {
		r.mu.Unlock()
		return errors.New("Order already exists")
	}

	r.orders[order.ID] = order
	r.mu.Unlock()
	return nil
}

func (r *OrderRepositoryLocal) DeleteOrder(id string) error {
	r.mu.Lock()

	if _, ok := r.orders[uuid.MustParse(id)]; !ok {
		r.mu.Unlock()
		return errors.New("Order not found")
	}

	delete(r.orders, uuid.MustParse(id))
	r.mu.Unlock()
	return nil
}

func (r *OrderRepositoryLocal) UpdateOrder(order *orderentity.Order) error {
	r.mu.Lock()
	r.orders[order.ID] = order
	r.mu.Unlock()
	return nil
}

func (r *OrderRepositoryLocal) GetOrderById(id string) (*orderentity.Order, error) {
	r.mu.Lock()

	if p, ok := r.orders[uuid.MustParse(id)]; ok {
		r.mu.Unlock()
		return p, nil
	}

	r.mu.Unlock()
	return nil, errors.New("Order not found")
}

func (r *OrderRepositoryLocal) GetOrderBy(key string, value string) (*orderentity.Order, error) {
	return nil, nil
}

func (r *OrderRepositoryLocal) GetAllOrder(key string, value string) ([]orderentity.Order, error) {
	orders := make([]orderentity.Order, 0)

	for _, p := range r.orders {
		orders = append(orders, *p)
	}

	return orders, nil
}
