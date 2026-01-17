package deliverydriverrepositorylocal

import (
	"context"
	"sync"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

// DeliveryDriverRepositoryLocal is an in-memory implementation of DeliveryDriverRepository
type DeliveryDriverRepositoryLocal struct {
	mu      sync.RWMutex
	drivers map[string]*model.DeliveryDriver
}

func NewDeliveryDriverRepositoryLocal() model.DeliveryDriverRepository {
	return &DeliveryDriverRepositoryLocal{drivers: make(map[string]*model.DeliveryDriver)}
}

func (r *DeliveryDriverRepositoryLocal) CreateDeliveryDriver(ctx context.Context, p *model.DeliveryDriver) error {

	r.drivers[p.ID.String()] = p
	return nil
}

func (r *DeliveryDriverRepositoryLocal) UpdateDeliveryDriver(ctx context.Context, p *model.DeliveryDriver) error {

	r.drivers[p.ID.String()] = p
	return nil
}

func (r *DeliveryDriverRepositoryLocal) DeleteDeliveryDriver(ctx context.Context, id string) error {

	delete(r.drivers, id)
	return nil
}

func (r *DeliveryDriverRepositoryLocal) GetDeliveryDriverById(ctx context.Context, id string) (*model.DeliveryDriver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if d, ok := r.drivers[id]; ok {
		return d, nil
	}
	return nil, nil
}

func (r *DeliveryDriverRepositoryLocal) GetDeliveryDriverByEmployeeId(ctx context.Context, id string) (*model.DeliveryDriver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, d := range r.drivers {
		if d.EmployeeID.String() == id {
			return d, nil
		}
	}
	return nil, nil
}

func (r *DeliveryDriverRepositoryLocal) GetAllDeliveryDrivers(ctx context.Context, isActive ...bool) ([]model.DeliveryDriver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.DeliveryDriver, 0, len(r.drivers))
	for _, d := range r.drivers {
		out = append(out, *d)
	}
	return out, nil
}
