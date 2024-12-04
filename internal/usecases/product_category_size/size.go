package productcategorysizeusecases

import (
	"context"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorysizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category_size"
)

type Service struct {
	rs productentity.SizeRepository
	rc productentity.CategoryRepository
}

func NewService(rs productentity.SizeRepository) *Service {
	return &Service{rs: rs}
}

func (s *Service) AddDependencies(rc productentity.CategoryRepository) {
	s.rc = rc
}

func (s *Service) CreateSize(ctx context.Context, dto *productcategorysizedto.CreateSizeInput) (uuid.UUID, error) {
	size, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	category, err := s.rc.GetCategoryById(ctx, size.CategoryID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if err = productentity.ValidateDuplicateSizes(size.Name, category.Sizes); err != nil {
		return uuid.Nil, err
	}

	err = s.rs.CreateSize(ctx, size)

	if err != nil {
		return uuid.Nil, err
	}

	return size.ID, nil
}

func (s *Service) UpdateSize(ctx context.Context, dtoId *entitydto.IdRequest, dto *productcategorysizedto.UpdateSizeInput) error {
	size, err := s.rs.GetSizeById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(size); err != nil {
		return err
	}

	category, err := s.rc.GetCategoryById(ctx, size.CategoryID.String())

	if err != nil {
		return err
	}

	if err = productentity.ValidateUpdateSize(size, category.Sizes); err != nil {
		return err
	}

	if err = s.rs.UpdateSize(ctx, size); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteSize(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.rs.GetSizeById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rs.DeleteSize(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}
func (s *Service) AddSizesByValues(ctx context.Context, dto *productcategorysizedto.CreateSizes) error {
	sizes, categoryID, err := dto.ToModel()
	if err != nil {
		return err
	}

	for _, size := range sizes {
		newSize := productentity.NewSize(productentity.SizeCommonAttributes{
			Name:       size,
			CategoryID: *categoryID,
		})

		if err := s.rs.CreateSize(ctx, newSize); err != nil {
			return err
		}
	}

	return nil
}
