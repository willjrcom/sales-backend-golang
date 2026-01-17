package clientrepositorylocal

import (
	"context"

	"errors"
	"sort"
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

// GetAllClients retrieves a paginated list of clients and the total count.
func (r *ClientRepositoryLocal) GetAllClients(ctx context.Context, page, perPage int, isActive ...bool) ([]model.Client, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.clients)
	if total == 0 {
		return []model.Client{}, 0, nil
	}
	ids := make([]string, 0, total)
	for id := range r.clients {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if page < 0 {
		page = 0
	}
	if perPage < 1 {
		perPage = total
	}
	offset := page * perPage
	if offset >= total {
		return []model.Client{}, total, nil
	}
	end := offset + perPage
	if end > total {
		end = total
	}
	segment := ids[offset:end]
	result := make([]model.Client, 0, len(segment))
	for _, id := range segment {
		result = append(result, *r.clients[id])
	}
	return result, total, nil
}
