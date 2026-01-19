package productcategoryusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
)

var (
	ErrSizeIsUsed     = errors.New("size is used in products")
	ErrProductsExists = errors.New("product category has products")
)

type Service struct {
	r  model.CategoryRepository
	sq quantityusecases.Service
	ss sizeusecases.Service
}

func NewService(c model.CategoryRepository, sq *quantityusecases.Service, ss *sizeusecases.Service) *Service {
	return &Service{r: c, sq: *sq, ss: *ss}
}

func (s *Service) CreateCategory(ctx context.Context, dto *productcategorydto.CategoryCreateDTO) (uuid.UUID, error) {
	category, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	if c, _ := s.r.GetCategoryByName(ctx, category.Name, false); c != nil {
		return uuid.Nil, errors.New("category name already exists")
	}

	categoryModel := &model.ProductCategory{}
	categoryModel.FromDomain(category)
	err = s.r.CreateCategory(ctx, categoryModel)

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
	categoryModel, err := s.r.GetCategoryById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	category := categoryModel.ToDomain()
	if err = dto.UpdateDomain(category); err != nil {
		return err
	}

	if c, _ := s.r.GetCategoryByName(ctx, category.Name, false); c != nil && c.ID != category.ID {
		return errors.New("category name already exists")
	}

	categoryModel.FromDomain(category)
	if err = s.r.UpdateCategory(ctx, categoryModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteCategoryById(ctx context.Context, dto *entitydto.IDRequest) error {
	category, err := s.r.GetCategoryById(ctx, dto.ID.String())
	if err != nil {
		return err
	}

	if len(category.Products) > 0 {
		return ErrProductsExists
	}

	if err := s.r.DeleteCategory(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetCategoryById(ctx context.Context, dto *entitydto.IDRequest) (*productcategorydto.CategoryDTO, error) {
	if categoryModel, err := s.r.GetCategoryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		category := categoryModel.ToDomain()
		dto := &productcategorydto.CategoryDTO{}
		dto.FromDomain(category)
		return dto, nil
	}
}

func (s *Service) GetAllCategories(ctx context.Context, dto *entitydto.IDRequest, page int, perPage int, isActive bool) ([]productcategorydto.CategoryDTO, error) {
	if categoryModels, err := s.r.GetAllCategories(ctx, dto.IDs, page, perPage, isActive); err != nil {
		return nil, err
	} else {
		dtos := modelsToDTOs(categoryModels)
		return dtos, nil
	}
}

func (s *Service) GetAllCategoriesWithProcessRulesAndOrderProcess(ctx context.Context) ([]productcategorydto.CategoryWithOrderProcessDTO, error) {
	if categoryModels, err := s.r.GetAllCategoriesWithProcessRulesAndOrderProcess(ctx); err != nil {
		return nil, err
	} else {
		dtos := modelsWithOrderProcessToDTOs(categoryModels)
		return dtos, nil
	}
}

func modelsToDTOs(categoryModels []model.ProductCategory) []productcategorydto.CategoryDTO {
	DTOs := []productcategorydto.CategoryDTO{}

	for _, categoryModel := range categoryModels {
		category := categoryModel.ToDomain()
		c := &productcategorydto.CategoryDTO{}
		c.FromDomain(category)
		DTOs = append(DTOs, *c)
	}

	return DTOs
}

func modelsWithOrderProcessToDTOs(categoryModels []model.ProductCategoryWithOrderProcess) []productcategorydto.CategoryWithOrderProcessDTO {
	DTOs := []productcategorydto.CategoryWithOrderProcessDTO{}

	for _, categoryModel := range categoryModels {
		category := categoryModel.ToDomain()
		c := &productcategorydto.CategoryWithOrderProcessDTO{}
		c.FromDomain(category)
		DTOs = append(DTOs, *c)
	}

	return DTOs
}

func (s *Service) GetComplementProducts(ctx context.Context, categoryID string) ([]productcategorydto.ProductDTO, error) {
	products, err := s.r.GetComplementProducts(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return productsToDTOs(products), nil
}

func (s *Service) GetAdditionalProducts(ctx context.Context, categoryID string) ([]productcategorydto.ProductDTO, error) {
	products, err := s.r.GetAdditionalProducts(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return productsToDTOs(products), nil
}

func productsToDTOs(products []model.Product) []productcategorydto.ProductDTO {
	dtos := []productcategorydto.ProductDTO{}
	for _, p := range products {
		domain := p.ToDomain()
		dto := &productcategorydto.ProductDTO{}
		dto.FromDomain(domain)
		dtos = append(dtos, *dto)
	}
	return dtos
}

func (s *Service) GetAllCategoriesMap(ctx context.Context, isActive bool) ([]productcategorydto.CategoryMapDTO, error) {
	if categoryModels, err := s.r.GetAllCategoriesMap(ctx, isActive); err != nil {
		return nil, err
	} else {
		dtos := make([]productcategorydto.CategoryMapDTO, 0)
		for _, model := range categoryModels {
			dtos = append(dtos, productcategorydto.CategoryMapDTO{
				ID:   model.ID,
				Name: model.Name,
			})
		}
		return dtos, nil
	}
}

func (s *Service) GetComplementCategories(ctx context.Context) ([]productcategorydto.CategoryMapDTO, error) {
	if categoryModels, err := s.r.GetComplementCategories(ctx); err != nil {
		return nil, err
	} else {
		dtos := make([]productcategorydto.CategoryMapDTO, 0)
		for _, model := range categoryModels {
			dtos = append(dtos, productcategorydto.CategoryMapDTO{
				ID:   model.ID,
				Name: model.Name,
			})
		}
		return dtos, nil
	}
}

func (s *Service) GetAdditionalCategories(ctx context.Context) ([]productcategorydto.CategoryMapDTO, error) {
	if categoryModels, err := s.r.GetAdditionalCategories(ctx); err != nil {
		return nil, err
	} else {
		dtos := make([]productcategorydto.CategoryMapDTO, 0)
		for _, model := range categoryModels {
			dtos = append(dtos, productcategorydto.CategoryMapDTO{
				ID:   model.ID,
				Name: model.Name,
			})
		}
		return dtos, nil
	}
}
