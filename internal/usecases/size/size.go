package sizeusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrProductsExists = errors.New("size is used in products")
)

type Service struct {
	rs model.SizeRepository
	rc model.CategoryRepository
}

func NewService(rs model.SizeRepository, rc model.CategoryRepository) *Service {
	return &Service{rs: rs, rc: rc}
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
	size, err := s.rs.GetSizeByIdWithProducts(ctx, dto.ID.String())
	if err != nil {
		return err
	}

	if len(size.Products) > 0 {
		return ErrProductsExists
	}

	if err := s.rs.DeleteSize(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetSizeById(ctx context.Context, dto *entitydto.IDRequest) (*sizedto.SizeDTO, error) {
	if sizeModel, err := s.rs.GetSizeById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		size := sizeModel.ToDomain()
		sizeDto := &sizedto.SizeDTO{}
		sizeDto.FromDomain(size)
		return sizeDto, nil
	}
}

func (s *Service) AddSizesByValues(ctx context.Context, dto *sizedto.SizeCreateBatchDTO) error {
	sizes, categoryID, err := dto.ToDomain()
	if err != nil {
		return err
	}

	for _, size := range sizes {
		newSize := productentity.NewSize(productentity.SizeCommonAttributes{
			Name:       size,
			IsActive:   true,
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

func (s *Service) GetSizesByCategoryId(ctx context.Context, categoryId string) ([]*sizedto.SizeDTO, error) {
	if sizes, err := s.rs.GetSizesByCategoryId(ctx, categoryId); err != nil {
		return nil, err
	} else {
		sizeDtos := []*sizedto.SizeDTO{}
		for _, size := range sizes {
			sDto := &sizedto.SizeDTO{}
			sDto.FromDomain(size.ToDomain())
			sizeDtos = append(sizeDtos, sDto)
		}
		return sizeDtos, nil
	}
}
