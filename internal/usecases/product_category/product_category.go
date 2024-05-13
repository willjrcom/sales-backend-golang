package productcategoryusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
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

func (s *Service) CreateCategory(ctx context.Context, dto *productcategorydto.CreateCategoryInput) (uuid.UUID, error) {
	category, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if c, _ := s.r.GetCategoryByName(ctx, category.Name, false); c != nil {
		return uuid.Nil, errors.New("category name already exists")
	}

	err = s.r.CreateCategory(ctx, category)

	if err != nil {
		return uuid.Nil, err
	}

	return category.ID, nil
}

func (s *Service) UpdateCategory(ctx context.Context, dtoId *entitydto.IdRequest, dto *productcategorydto.UpdateCategoryInput) error {
	category, err := s.r.GetCategoryById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(category); err != nil {
		return err
	}

	if c, _ := s.r.GetCategoryByName(ctx, category.Name, false); c != nil && c.ID != category.ID {
		return errors.New("category name already exists")
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

func (s *Service) GetCategoryById(ctx context.Context, dto *entitydto.IdRequest) (*productcategorydto.CategoryOutput, error) {
	if productCategory, err := s.r.GetCategoryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &productcategorydto.CategoryOutput{}
		dto.FromModel(productCategory)
		return dto, nil
	}
}

func (s *Service) GetAllCategories(ctx context.Context) ([]productcategorydto.CategoryOutput, error) {
	if categories, err := s.r.GetAllCategories(ctx); err != nil {
		return nil, err
	} else {
		dtos := categorySizesToDtos(categories)
		return dtos, nil
	}
}

func categorySizesToDtos(categories []productentity.ProductCategory) []productcategorydto.CategoryOutput {
	dtoOutput := []productcategorydto.CategoryOutput{}

	for _, category := range categories {
		c := &productcategorydto.CategoryOutput{}
		c.FromModel(&category)
		dtoOutput = append(dtoOutput, *c)
	}

	return dtoOutput
}
