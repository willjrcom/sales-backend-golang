package categoryproductusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	categorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/category"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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

func (s *Service) RegisterCategory(ctx context.Context, dto *categorydto.RegisterCategoryInput) (uuid.UUID, error) {
	category, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if c, _ := s.r.GetCategoryByName(ctx, category.Name, false); c != nil {
		return uuid.Nil, errors.New("category name already exists")
	}

	err = s.r.RegisterCategory(ctx, category)

	if err != nil {
		return uuid.Nil, err
	}

	return category.ID, nil
}

func (s *Service) UpdateCategory(ctx context.Context, dtoId *entitydto.IdRequest, dto *categorydto.UpdateCategoryInput) error {
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

func (s *Service) GetCategoryById(ctx context.Context, dto *entitydto.IdRequest) (*categorydto.CategoryOutput, error) {
	if categoryProduct, err := s.r.GetCategoryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &categorydto.CategoryOutput{}
		dto.FromModel(categoryProduct)
		return dto, nil
	}
}

func (s *Service) GetAllCategories(ctx context.Context) ([]categorydto.CategoryOutput, error) {
	if categories, err := s.r.GetAllCategories(ctx); err != nil {
		return nil, err
	} else {
		dtos := categorySizesToDtos(categories)
		return dtos, nil
	}
}

func categorySizesToDtos(categories []productentity.Category) []categorydto.CategoryOutput {
	dtoOutput := []categorydto.CategoryOutput{}

	for _, category := range categories {
		c := &categorydto.CategoryOutput{}
		c.FromModel(&category)
		dtoOutput = append(dtoOutput, *c)
	}

	return dtoOutput
}
