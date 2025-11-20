package orderrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderRepositoryLocal struct {
	orders map[uuid.UUID]*model.Order
	mu     sync.RWMutex
}

func NewOrderRepositoryLocal() model.OrderRepository {
	return &OrderRepositoryLocal{
		orders: make(map[uuid.UUID]*model.Order),
	}
}

func (r *OrderRepositoryLocal) CreateOrder(ctx context.Context, order *model.Order) error {
	if order == nil || order.Entity.ID == uuid.Nil {
		return errors.New("invalid order")
	}

	if _, exists := r.orders[order.Entity.ID]; exists {
		return errors.New("order already exists")
	}
	r.orders[order.Entity.ID] = order
	return nil
}

func (r *OrderRepositoryLocal) DeleteOrder(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid id")
	}
	urid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if _, exists := r.orders[urid]; !exists {
		return errors.New("order not found")
	}
	delete(r.orders, urid)
	return nil
}

func (r *OrderRepositoryLocal) PendingOrder(ctx context.Context, order *model.Order) error {
	if order == nil || order.Entity.ID == uuid.Nil {
		return errors.New("invalid order")
	}

	r.orders[order.Entity.ID] = order
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

func (r *OrderRepositoryLocal) GetAllOrders(ctx context.Context, shiftID string, withStatus []orderentity.StatusOrder, withCategory bool, queryCondition string) ([]model.Order, error) {
	orders := make([]model.Order, 0)

	for _, p := range r.orders {
		orders = append(orders, *p)
	}

	return orders, nil
}

// GetAllOrdersWithDelivery returns orders with delivery information, paginated by page and perPage
func (r *OrderRepositoryLocal) GetAllOrdersWithDelivery(ctx context.Context, shiftID string, page, perPage int) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := []model.Order{}
	for _, o := range r.orders {
		if o.Delivery != nil {
			all = append(all, *o)
		}
	}
	if page < 1 || perPage < 1 {
		return []model.Order{}, nil
	}
	start := (page - 1) * perPage
	if start >= len(all) {
		return []model.Order{}, nil
	}
	end := start + perPage
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], nil
}

// GetAllOrdersWithDelivery returns orders with delivery information, paginated by page and perPage
func (r *OrderRepositoryLocal) GetAllOrdersWithPickup(ctx context.Context, shiftID string, page, perPage int) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := []model.Order{}
	for _, o := range r.orders {
		if o.Delivery != nil {
			all = append(all, *o)
		}
	}
	if page < 1 || perPage < 1 {
		return []model.Order{}, nil
	}
	start := (page - 1) * perPage
	if start >= len(all) {
		return []model.Order{}, nil
	}
	end := start + perPage
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], nil
}

func (r *OrderRepositoryLocal) UpdateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {

	orderID := delivery.OrderID.String()
	for _, o := range r.orders {
		if o.ID.String() == orderID {
			o.Delivery = delivery
			return nil
		}
	}
	return errors.New("order not found")
}

func (r *OrderRepositoryLocal) DeleteOrderDelivery(ctx context.Context, id string) error {

	for _, o := range r.orders {
		if o.Delivery != nil && o.Delivery.ID.String() == id {
			o.Delivery = nil
			return nil
		}
	}
	return errors.New("delivery not found")
}

func (r *OrderRepositoryLocal) UpdateOrderTable(ctx context.Context, table *model.OrderTable) error {

	orderID := table.OrderID.String()
	for _, o := range r.orders {
		if o.ID.String() == orderID {
			o.Table = table
			return nil
		}
	}
	return errors.New("order not found")
}

func (r *OrderRepositoryLocal) DeleteOrderTable(ctx context.Context, id string) error {

	for _, o := range r.orders {
		if o.Table != nil && o.Table.ID.String() == id {
			o.Table = nil
			return nil
		}
	}
	return errors.New("table not found")
}

func (r *OrderRepositoryLocal) AddPaymentOrder(ctx context.Context, payment *model.PaymentOrder) error {

	orderID := payment.OrderID.String()
	for _, o := range r.orders {
		if o.ID.String() == orderID {
			o.Payments = append(o.Payments, *payment)
			return nil
		}
	}
	return errors.New("order not found")
}
