package sizerepositorylocal

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	errSizeExists   = errors.New("size already exists")
	errSizeNotFound = errors.New("size not found")
)

type SizeRepositoryLocal struct {
	sizes map[uuid.UUID]*model.Size
}

func NewSizeRepositoryLocal() model.SizeRepository {
	return &SizeRepositoryLocal{sizes: make(map[uuid.UUID]*model.Size)}
}

func (r *SizeRepositoryLocal) CreateSize(_ context.Context, p *model.Size) error {

	if _, ok := r.sizes[p.ID]; ok {
		return errSizeExists
	}

	r.sizes[p.ID] = p
	return nil
}

func (r *SizeRepositoryLocal) UpdateSize(_ context.Context, s *model.Size) error {
	r.sizes[s.ID] = s
	return nil
}

func (r *SizeRepositoryLocal) DeleteSize(_ context.Context, id string) error {

	if _, ok := r.sizes[uuid.MustParse(id)]; !ok {
		return errSizeNotFound
	}

	delete(r.sizes, uuid.MustParse(id))
	return nil
}

func (r *SizeRepositoryLocal) GetSizeById(_ context.Context, id string) (*model.Size, error) {

	if p, ok := r.sizes[uuid.MustParse(id)]; ok {
		return p, nil
	}

	return nil, errSizeNotFound
}
