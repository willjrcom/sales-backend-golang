package productusecases

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
)

type Service struct {
	Repository productentity.Repository
}

func NewService(repository productentity.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) RegisterProduct(dto *productdto.CreateProductInput) (uuid.UUID, error) {
	product, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err := s.Repository.RegisterProduct(product); err != nil {
		return uuid.Nil, err
	}

	return product.ID, nil
}

func (s *Service) UpdateProduct(dtoId entitydto.IdRequest, dto *productdto.UpdateProductInput) error {
	product, err := s.Repository.GetProduct(dtoId.Id.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(product); err != nil {
		return err
	}

	if err := s.Repository.UpdateProduct(product); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProduct(dto *entitydto.IdRequest) error {
	if _, err := s.Repository.GetProduct(dto.Id.String()); err != nil {
		return err
	}

	if err := s.Repository.DeleteProduct(dto.Id.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProduct(dto *entitydto.IdRequest) (*productentity.Product, error) {
	if product, err := s.Repository.GetProduct(dto.Id.String()); err != nil {
		return nil, err
	} else {
		return product, nil
	}
}
