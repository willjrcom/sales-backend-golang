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
	r productentity.CategoryRepository
}

func NewService(c productentity.CategoryRepository) *Service {
	return &Service{r: c}
}

func (s *Service) RegisterCategory(ctx context.Context, dto *productdto.RegisterCategoryInput) (uuid.UUID, error) {
	categoryProduct, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.RegisterCategory(ctx, categoryProduct)

	if err != nil {
		return uuid.Nil, err
	}

	return categoryProduct.ID, nil
}

func (s *Service) UpdateCategory(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateCategoryInput) error {
	category, err := s.r.GetCategoryById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(category); err != nil {
		return err
	}

	if err = s.r.UpdateCategory(ctx, category); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteCategoryById(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetCategoryById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteCategory(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetCategoryById(ctx context.Context, dto *entitydto.IdRequest) (*productdto.CategorySizesOutput, error) {
	if categoryProduct, err := s.r.GetCategoryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &productdto.CategorySizesOutput{}
		dto.FromModel(categoryProduct)
		return dto, nil
	}
}

func (s *Service) GetAllCategoryProducts(ctx context.Context) ([]productdto.CategoryProductsOutput, error) {
	if categories, err := s.r.GetAllCategoryProducts(ctx); err != nil {
		return nil, err
	} else {
		dtos := categoryProductsToDtos(categories)
		return dtos, nil
	}
}

func (s *Service) GetAllCategorySizes(ctx context.Context) ([]productdto.CategorySizesOutput, error) {
	if categories, err := s.r.GetAllCategorySizes(ctx); err != nil {
		return nil, err
	} else {
		dtos := categorySizesToDtos(categories)
		return dtos, nil
	}
}

func categoryProductsToDtos(categories []productentity.Category) []productdto.CategoryProductsOutput {
	dtoOutput := []productdto.CategoryProductsOutput{}

	for _, category := range categories {
		c := &productdto.CategoryProductsOutput{}
		c.FromModel(&category)
		dtoOutput = append(dtoOutput, *c)
	}

	return dtoOutput
}

func categorySizesToDtos(categories []productentity.Category) []productdto.CategorySizesOutput {
	dtoOutput := []productdto.CategorySizesOutput{}

	for _, category := range categories {
		c := &productdto.CategorySizesOutput{}
		c.FromModel(&category)
		dtoOutput = append(dtoOutput, *c)
	}

	return dtoOutput
}
