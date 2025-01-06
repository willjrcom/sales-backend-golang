package sizeusecases

import (
	"context"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	rs model.SizeRepository
	rc model.CategoryRepository
}

func NewService(rs model.SizeRepository) *Service {
	return &Service{rs: rs}
}

func (s *Service) AddDependencies(rc model.CategoryRepository) {
	s.rc = rc
}

func (s *Service) CreateSize(ctx context.Context, dto *sizedto.SizeCreateDTO) (uuid.UUID, error) {
	size, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	categoryModel, err := s.rc.GetCategoryById(ctx, size.CategoryID.String())

	if err != nil {
		return uuid.Nil, err
	}

	category := categoryModel.ToDomain()

	if err = productentity.ValidateDuplicateSizes(size.Name, category.Sizes); err != nil {
		return uuid.Nil, err
	}

	sizeModel := &model.Size{}
	sizeModel.FromDomain(size)
	err = s.rs.CreateSize(ctx, sizeModel)

	if err != nil {
		return uuid.Nil, err
	}

	return size.ID, nil
}

func (s *Service) UpdateSize(ctx context.Context, dtoId *entitydto.IDRequest, dto *sizedto.SizeUpdateDTO) error {
	sizeModel, err := s.rs.GetSizeById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	size := sizeModel.ToDomain()
	if err = dto.UpdateDomain(size); err != nil {
		return err
	}

	categoryModel, err := s.rc.GetCategoryById(ctx, size.CategoryID.String())

	if err != nil {
		return err
	}

	category := categoryModel.ToDomain()

	if err = productentity.ValidateUpdateSize(size, category.Sizes); err != nil {
		return err
	}

	sizeModel.FromDomain(size)
	if err = s.rs.UpdateSize(ctx, sizeModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteSize(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.rs.GetSizeById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rs.DeleteSize(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}
func (s *Service) AddSizesByValues(ctx context.Context, dto *sizedto.SizeCreateBatchDTO) error {
	sizes, categoryID, err := dto.ToDomain()
	if err != nil {
		return err
	}

	isTrue := true

	for _, size := range sizes {
		newSize := productentity.NewSize(productentity.SizeCommonAttributes{
			Name:       size,
			IsActive:   &isTrue,
			CategoryID: *categoryID,
		})

		sizeModel := &model.Size{}
		sizeModel.FromDomain(newSize)
		if err := s.rs.CreateSize(ctx, sizeModel); err != nil {
			return err
		}
	}

	return nil
}
