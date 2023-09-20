package sizeusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
)

var (
	ErrSizeIsUsed = errors.New("size is used in products")
)

type Service struct {
	r productentity.SizeRepository
}

func NewService(c productentity.SizeRepository) *Service {
	return &Service{r: c}
}

func (s *Service) RegisterSize(ctx context.Context, dto *productdto.RegisterSizeInput) (uuid.UUID, error) {
	size, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.RegisterSize(ctx, size)

	if err != nil {
		return uuid.Nil, err
	}

	return size.ID, nil
}

func (s *Service) UpdateSize(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateSizeInput) error {
	size, err := s.r.GetSizeById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(size); err != nil {
		return err
	}

	if err = s.r.UpdateSize(ctx, size); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteSize(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetSizeById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteSize(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetSizeById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.Size, error) {
	if size, err := s.r.GetSizeById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return size, nil
	}
}
