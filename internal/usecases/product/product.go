package productusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
)

var (
	ErrSizeIsInvalid = errors.New("size is invalid")
)

type Service struct {
	rProduct  productentity.Repository
	rCategory productentity.RepositoryCategory
}

func NewService(r productentity.Repository, c productentity.RepositoryCategory) *Service {
	return &Service{rProduct: r, rCategory: c}
}

func (s *Service) RegisterProduct(ctx context.Context, dto *productdto.RegisterProductInput) (uuid.UUID, error) {
	product, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if category, err := s.rCategory.GetCategoryProductById(ctx, product.CategoryID.String()); err == nil {
		product.Category = category
	} else {
		return uuid.Nil, err
	}

	if exists, err := product.FindSizeInCategory(); !exists {
		return uuid.Nil, err
	}

	if err := s.rProduct.RegisterProduct(ctx, product); err != nil {
		return uuid.Nil, err
	}

	return product.ID, nil
}

func (s *Service) UpdateProduct(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateProductInput) error {
	product, err := s.rProduct.GetProductById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(product); err != nil {
		return err
	}

	if err := s.rProduct.UpdateProduct(ctx, product); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProductById(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.rProduct.GetProductById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rProduct.DeleteProduct(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProductById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.Product, error) {
	if product, err := s.rProduct.GetProductById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return product, nil
	}
}

func (s *Service) GetProductBy(ctx context.Context, filter *productdto.FilterProductInput) ([]productentity.Product, error) {
	product := filter.ToModel()

	if products, err := s.rProduct.GetProductsBy(ctx, product); err != nil {
		return nil, err
	} else {
		return products, nil
	}
}

func (s *Service) GetAllProductsByCategory(ctx context.Context, dto *filterdto.Category) ([]productentity.Product, error) {
	products, err := s.rProduct.GetAllProductsByCategory(ctx, dto.Category)

	if err != nil {
		return nil, err
	}

	return products, nil
}
