package categoryproductusecases

import (
	"context"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
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

func (s *Service) UpdateCategoryProduct(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateCategoryProductInput) error {
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
