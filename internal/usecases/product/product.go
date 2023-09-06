package productusecases

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
)

type Service struct {
	rProduct  productentity.Repository
	rCategory productentity.RepositoryCategory
}

func NewService(r productentity.Repository, c productentity.RepositoryCategory) *Service {
	return &Service{rProduct: r, rCategory: c}
}

func (s *Service) RegisterProduct(dto *productdto.CreateProductInput) (uuid.UUID, error) {
	product, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err := s.rProduct.RegisterProduct(product); err != nil {
		return uuid.Nil, err
	}

	return product.ID, nil
}

func (s *Service) UpdateProduct(dtoId *entitydto.IdRequest, dto *productdto.UpdateProductInput) error {
	product, err := s.rProduct.GetProductById(dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(product); err != nil {
		return err
	}

	if err := s.rProduct.UpdateProduct(product); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProductById(dto *entitydto.IdRequest) error {
	if _, err := s.rProduct.GetProductById(dto.ID.String()); err != nil {
		return err
	}

	if err := s.rProduct.DeleteProduct(dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProductById(dto *entitydto.IdRequest) (*productentity.Product, error) {
	if product, err := s.rProduct.GetProductById(dto.ID.String()); err != nil {
		return nil, err
	} else {
		return product, nil
	}
}

func (s *Service) GetAllProduct(dto *filterdto.Filter) ([]productentity.Product, error) {
	if products, err := s.rProduct.GetAllProduct(dto.Key, dto.Value); err != nil {
		return nil, err
	} else {
		return products, nil
	}
}
