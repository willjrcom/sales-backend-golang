package contactrepositorylocal

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ContactRepositoryLocal struct {
	contacts map[string]*model.Contact
	mu       sync.RWMutex
}

func NewContactRepositoryLocal() model.ContactRepository {
	return &ContactRepositoryLocal{
		contacts: make(map[string]*model.Contact),
	}
}

func (r *ContactRepositoryLocal) CreateContact(ctx context.Context, c *model.Contact) error {
	if c == nil || c.ID == uuid.Nil {
		return errors.New("invalid contact")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.contacts[c.ID.String()]; exists {
		return errors.New("contact already exists")
	}
	r.contacts[c.ID.String()] = c
	return nil
}

func (r *ContactRepositoryLocal) UpdateContact(ctx context.Context, c *model.Contact) error {
	if c == nil || c.ID == uuid.Nil {
		return errors.New("invalid contact")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.contacts[c.ID.String()]; !exists {
		return errors.New("contact not found")
	}
	r.contacts[c.ID.String()] = c
	return nil
}

func (r *ContactRepositoryLocal) DeleteContact(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid id")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.contacts[id]; !exists {
		return errors.New("contact not found")
	}
	delete(r.contacts, id)
	return nil
}

func (r *ContactRepositoryLocal) GetContactById(ctx context.Context, id string) (*model.Contact, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.contacts[id]
	if !exists {
		return nil, errors.New("contact not found")
	}
	return c, nil
}

func (r *ContactRepositoryLocal) GetContactByDddAndNumber(ctx context.Context, ddd string, number string, contactType string) (*model.Contact, error) {
	if ddd == "" || number == "" {
		return nil, errors.New("ddd and number are required")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.contacts {
		if c.Ddd == ddd && c.Number == number && (contactType == "" || strings.EqualFold(c.Type, contactType)) {
			return c, nil
		}
	}
	return nil, errors.New("contact not found")
}

func (r *ContactRepositoryLocal) FtSearchContacts(ctx context.Context, key string, contactType string) ([]model.Contact, error) {
	var results []model.Contact
	r.mu.RLock()
	defer r.mu.RUnlock()
	keyLower := strings.ToLower(key)
	for _, c := range r.contacts {
		if contactType != "" && !strings.EqualFold(c.Type, contactType) {
			continue
		}
		if strings.Contains(strings.ToLower(c.Number), keyLower) || strings.Contains(strings.ToLower(c.Ddd), keyLower) {
			results = append(results, *c)
		}
	}
	return results, nil
}
