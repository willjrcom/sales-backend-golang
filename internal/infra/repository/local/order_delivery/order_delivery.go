package orderdeliveryrepositorylocal

import (
	"context"
	"sync"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderDeliveryRepositoryLocal struct {
	mu         sync.RWMutex
	deliveries map[string]*model.OrderDelivery
}

func NewOrderDeliveryRepositoryLocal() model.OrderDeliveryRepository {
	return &OrderDeliveryRepositoryLocal{deliveries: make(map[string]*model.OrderDelivery)}
}

func (r *OrderDeliveryRepositoryLocal) CreateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deliveries[delivery.ID.String()] = delivery
	return nil
}

func (r *OrderDeliveryRepositoryLocal) UpdateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deliveries[delivery.ID.String()] = delivery
	return nil
}

func (r *OrderDeliveryRepositoryLocal) DeleteOrderDelivery(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.deliveries, id)
	return nil
}

func (r *OrderDeliveryRepositoryLocal) GetDeliveryById(ctx context.Context, id string) (*model.OrderDelivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if d, ok := r.deliveries[id]; ok {
		return d, nil
	}
	return nil, nil
}

func (r *OrderDeliveryRepositoryLocal) GetDeliveriesByIds(ctx context.Context, ids []string) ([]model.OrderDelivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []model.OrderDelivery{}
	for _, id := range ids {
		if d, ok := r.deliveries[id]; ok {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (r *OrderDeliveryRepositoryLocal) GetDeliveriesByClientId(ctx context.Context, clientID string) ([]model.OrderDelivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []model.OrderDelivery{}
	for _, d := range r.deliveries {
		if d.ClientID.String() == clientID {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (r *OrderDeliveryRepositoryLocal) GetAllDeliveries(ctx context.Context) ([]model.OrderDelivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.OrderDelivery, 0, len(r.deliveries))
	for _, d := range r.deliveries {
		out = append(out, *d)
	}
	return out, nil
}
