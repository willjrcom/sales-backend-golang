package sizerepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	errSizeExists   = errors.New("size already exists")
	errSizeNotFound = errors.New("size not found")
)

type SizeRepositoryLocal struct {
	mu    sync.Mutex
	sizes map[uuid.UUID]*productentity.Size
}

func NewSizeRepositoryLocal() *SizeRepositoryLocal {
	return &SizeRepositoryLocal{sizes: make(map[uuid.UUID]*productentity.Size)}
}

func (r *SizeRepositoryLocal) RegisterSize(_ context.Context, p *productentity.Size) error {
	r.mu.Lock()

	if _, ok := r.sizes[p.ID]; ok {
		r.mu.Unlock()
		return errSizeExists
	}

	r.sizes[p.ID] = p
	r.mu.Unlock()
	return nil
}

func (r *SizeRepositoryLocal) UpdateSize(_ context.Context, s *productentity.Size) error {
	r.mu.Lock()
	r.sizes[s.ID] = s
	r.mu.Unlock()
	return nil
}

func (r *SizeRepositoryLocal) DeleteSize(_ context.Context, id string) error {
	r.mu.Lock()

	if _, ok := r.sizes[uuid.MustParse(id)]; !ok {
		r.mu.Unlock()
		return errSizeNotFound
	}

	delete(r.sizes, uuid.MustParse(id))
	r.mu.Unlock()
	return nil
}

func (r *SizeRepositoryLocal) GetSizeById(_ context.Context, id string) (*productentity.Size, error) {
	r.mu.Lock()

	if p, ok := r.sizes[uuid.MustParse(id)]; ok {
		r.mu.Unlock()
		return p, nil
	}

	r.mu.Unlock()
	return nil, errSizeNotFound
}
