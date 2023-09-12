package categoryproductusecases

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
	r productentity.RepositoryCategory
}

func NewService(c productentity.RepositoryCategory) *Service {
	return &Service{r: c}
}

func (s *Service) RegisterCategoryProduct(ctx context.Context, dto *productdto.RegisterCategoryProductInput) (uuid.UUID, error) {
	categoryProduct, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.RegisterCategoryProduct(ctx, categoryProduct)

	if err != nil {
		return uuid.Nil, err
	}

	return categoryProduct.ID, nil
}

func (s *Service) UpdateCategoryProductName(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateCategoryProductNameInput) error {
	category, err := s.r.GetCategoryProductById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(category); err != nil {
		return err
	}

	if err = s.r.UpdateCategoryProduct(ctx, category); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateCategoryProductSizes(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateCategoryProductSizesInput) error {
	category, err := s.r.GetCategoryProductById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	// Validate if product use size
	categoryDb, err := s.r.GetCategoryProductById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}
	if exists, err := validateSizeToUpdate(ctx, categoryDb.Products); err != nil {
		return err
	} else if exists {
		return ErrSizeIsUsed
	}

	if err = dto.UpdateModel(category); err != nil {
		return err
	}

	if err = s.r.UpdateCategoryProduct(ctx, category); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteCategoryProductById(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetCategoryProductById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteCategoryProduct(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetCategoryProductById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.CategoryProduct, error) {
	if categoryProduct, err := s.r.GetCategoryProductById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return categoryProduct, nil
	}
}

func (s *Service) GetAllCategoryProduct(ctx context.Context) ([]productentity.CategoryProduct, error) {
	if categories, err := s.r.GetAllCategoryProduct(ctx); err != nil {
		return nil, err
	} else {
		return categories, nil
	}
}

func validateSizeToUpdate(ctx context.Context, products []productentity.Product) (bool, error) {
	for _, product := range products {
		if exists, err := product.FindSizeInCategory(); err != nil {
			return false, err
		} else if exists {
			return exists, nil
		}
	}

	return false, nil
}
