package productusecases

import (
	"github.com/google/uuid"
	productEntity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) RegisterProduct(dto *productdto.CreateProductInput) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (s *Service) UpdateProduct(dto *productdto.UpdateProductInput) error {
	return nil
}

func (s *Service) DeleteProduct(dto *entitydto.IdRequest) error {
	return nil
}

func (s *Service) GetProduct(dto *entitydto.IdRequest) (*productEntity.Product, error) {
	return nil, nil
}
