package productcategoryusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
)

var (
	ErrSizeIsUsed = errors.New("size is used in products")
)

type Service struct {
	r  productentity.CategoryRepository
	sq quantityusecases.Service
	ss sizeusecases.Service
}

func NewService(c productentity.CategoryRepository) *Service {
	return &Service{r: c}
}

func (s *Service) AddDependencies(sq quantityusecases.Service, ss sizeusecases.Service) {
	s.ss = ss
	s.sq = sq
}

func (s *Service) CreateCategory(ctx context.Context, dto *productcategorydto.CategoryCreateDTO) (uuid.UUID, error) {
	category, err := dto.ToDomain()

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

	quantities := []float64{0.3, 0.4, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sizes := []string{"P", "M", "G"}

	if category.IsAdditional {
		quantities = []float64{1, 2, 3, 4, 5}
		sizes = []string{"Padrão"}
	}

	registerQuantities := &quantitydto.QuantityCreateBatchDTO{
		Quantities: quantities,
		CategoryID: category.ID,
	}

	if err := s.sq.AddQuantitiesByValues(ctx, registerQuantities); err != nil {
		return category.ID, err
	}

	registerSizes := &sizedto.SizeCreateBatchDTO{
		Sizes:      sizes,
		CategoryID: category.ID,
	}

	if err := s.ss.AddSizesByValues(ctx, registerSizes); err != nil {
		return category.ID, err
	}
	return category.ID, nil
}

func (s *Service) UpdateCategory(ctx context.Context, dtoId *entitydto.IDRequest, dto *productcategorydto.CategoryUpdateDTO) error {
	category, err := s.r.GetCategoryById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateDomain(category); err != nil {
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

func (s *Service) DeleteCategoryById(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.r.GetCategoryById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteCategory(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetCategoryById(ctx context.Context, dto *entitydto.IDRequest) (*productcategorydto.CategoryDTO, error) {
	if productCategory, err := s.r.GetCategoryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &productcategorydto.CategoryDTO{}
		dto.FromDomain(productCategory)
		return dto, nil
	}
}

func (s *Service) GetAllCategories(ctx context.Context) ([]productcategorydto.CategoryDTO, error) {
	if categories, err := s.r.GetAllCategories(ctx); err != nil {
		return nil, err
	} else {
		dtos := categorySizesToDtos(categories)
		return dtos, nil
	}
}

func categorySizesToDtos(categories []productentity.ProductCategory) []productcategorydto.CategoryDTO {
	DTOs := []productcategorydto.CategoryDTO{}

	for _, category := range categories {
		c := &productcategorydto.CategoryDTO{}
		c.FromDomain(&category)
		DTOs = append(DTOs, *c)
	}

	return DTOs
}
