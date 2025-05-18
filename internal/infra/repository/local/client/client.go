package clientrepositorylocal

import (
	"context"

	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ClientRepositoryLocal struct {
	clients map[string]*model.Client
	mu      sync.RWMutex
}

func NewClientRepositoryLocal() model.ClientRepository {
	return &ClientRepositoryLocal{
		clients: make(map[string]*model.Client),
	}
}

func (r *ClientRepositoryLocal) CreateClient(ctx context.Context, p *model.Client) error {
	if p == nil || p.ID == uuid.Nil {
		return errors.New("invalid client")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.clients[p.ID.String()]; exists {
		return errors.New("client already exists")
	}
	r.clients[p.ID.String()] = p
	return nil
}

func (r *ClientRepositoryLocal) UpdateClient(ctx context.Context, p *model.Client) error {
	if p == nil || p.ID == uuid.Nil {
		return errors.New("invalid client")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.clients[p.ID.String()]; !exists {
		return errors.New("client not found")
	}
	r.clients[p.ID.String()] = p
	return nil
}

func (r *ClientRepositoryLocal) DeleteClient(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid id")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.clients[id]; !exists {
		return errors.New("client not found")
	}
	delete(r.clients, id)
	return nil
}

func (r *ClientRepositoryLocal) GetClientById(ctx context.Context, id string) (*model.Client, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.clients[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return c, nil
}

func (r *ClientRepositoryLocal) GetAllClients(ctx context.Context) ([]model.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var clients []model.Client
	for _, c := range r.clients {
		clients = append(clients, *c)
	}
	return clients, nil
}
