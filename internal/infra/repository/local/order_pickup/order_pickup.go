package orderpickuprepositorylocal

import (
	"context"
	"sync"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderPickupRepositoryLocal struct {
	mu      sync.RWMutex
	pickups map[string]*model.OrderPickup
}

func NewOrderPickupRepositoryLocal() model.OrderPickupRepository {
	return &OrderPickupRepositoryLocal{pickups: make(map[string]*model.OrderPickup)}
}

func (r *OrderPickupRepositoryLocal) CreateOrderPickup(ctx context.Context, pickup *model.OrderPickup) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pickups[pickup.ID.String()] = pickup
	return nil
}

func (r *OrderPickupRepositoryLocal) UpdateOrderPickup(ctx context.Context, pickup *model.OrderPickup) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pickups[pickup.ID.String()] = pickup
	return nil
}

func (r *OrderPickupRepositoryLocal) DeleteOrderPickup(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.pickups, id)
	return nil
}

func (r *OrderPickupRepositoryLocal) GetPickupById(ctx context.Context, id string) (*model.OrderPickup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.pickups[id]; ok {
		return p, nil
	}
	return nil, nil
}

func (r *OrderPickupRepositoryLocal) GetPickupsByContact(ctx context.Context, contact string) ([]model.OrderPickup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []model.OrderPickup
	for _, p := range r.pickups {
		if p.Contact == contact {
			out = append(out, *p)
		}
	}
	return out, nil
}

func (r *OrderPickupRepositoryLocal) GetAllPickups(ctx context.Context) ([]model.OrderPickup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.OrderPickup, 0, len(r.pickups))
	for _, p := range r.pickups {
		out = append(out, *p)
	}
	return out, nil
}
