package addressrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type AddressRepositoryLocal struct {
	addresses map[string]*model.Address
	mu        sync.RWMutex
}

func NewAddressRepositoryLocal() model.AddressRepository {
	return &AddressRepositoryLocal{
		addresses: make(map[string]*model.Address),
	}
}

func (r *AddressRepositoryLocal) CreateAddress(ctx context.Context, address *model.Address) error {
	if address == nil || address.ID == uuid.Nil {
		return errors.New("invalid address")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.addresses[address.ID.String()]; exists {
		return errors.New("address already exists")
	}
	r.addresses[address.ID.String()] = address
	return nil
}

func (r *AddressRepositoryLocal) UpdateAddress(ctx context.Context, address *model.Address) error {
	if address == nil || address.ID == uuid.Nil {
		return errors.New("invalid address")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.addresses[address.ID.String()]; !exists {
		return errors.New("address not found")
	}
	r.addresses[address.ID.String()] = address
	return nil
}

func (r *AddressRepositoryLocal) DeleteAddress(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid id")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.addresses[id]; !exists {
		return errors.New("address not found")
	}
	delete(r.addresses, id)
	return nil
}

func (r *AddressRepositoryLocal) GetAddressById(ctx context.Context, id string) (*model.Address, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, exists := r.addresses[id]
	if !exists {
		return nil, errors.New("address not found")
	}
	return a, nil
}
