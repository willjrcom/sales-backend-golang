package quantityusecases

import (
	"context"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	rq model.QuantityRepository
	rc model.CategoryRepository
}

func NewService(rq model.QuantityRepository, rc model.CategoryRepository) *Service {
	return &Service{rq: rq, rc: rc}
}

func (s *Service) CreateQuantity(ctx context.Context, dto *quantitydto.QuantityCreateDTO) (uuid.UUID, error) {
	quantity, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	categoryModel, err := s.rc.GetCategoryById(ctx, quantity.CategoryID.String())

	if err != nil {
		return uuid.Nil, err
	}

	category := categoryModel.ToDomain()
	if err = productentity.ValidateDuplicateQuantities(quantity.Quantity, category.Quantities); err != nil {
		return uuid.Nil, err
	}

	quantityModel := &model.Quantity{}
	quantityModel.FromDomain(quantity)
	if err = s.rq.CreateQuantity(ctx, quantityModel); err != nil {
		return uuid.Nil, err
	}

	return quantity.ID, nil
}

func (s *Service) UpdateQuantity(ctx context.Context, dtoId *entitydto.IDRequest, dto *quantitydto.QuantityUpdateDTO) error {
	quantityModel, err := s.rq.GetQuantityById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	quantity := quantityModel.ToDomain()
	if err = dto.UpdateDomain(quantity); err != nil {
		return err
	}

	categoryModel, err := s.rc.GetCategoryById(ctx, quantity.CategoryID.String())

	if err != nil {
		return err
	}

	category := categoryModel.ToDomain()

	if err = productentity.ValidateUpdateQuantity(quantity, category.Quantities); err != nil {
		return err
	}

	quantityModel.FromDomain(quantity)
	if err = s.rq.UpdateQuantity(ctx, quantityModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteQuantity(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.rq.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rq.DeleteQuantity(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQuantityById(ctx context.Context, dto *entitydto.IDRequest) (*quantitydto.QuantityDTO, error) {
	if quantityModel, err := s.rq.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		quantity := quantityModel.ToDomain()

		quantityDto := &quantitydto.QuantityDTO{}
		quantityDto.FromDomain(quantity)
		return quantityDto, nil
	}
}

func (s *Service) AddQuantitiesByValues(ctx context.Context, dto *quantitydto.QuantityCreateBatchDTO) error {
	quantities, categoryID, err := dto.ToDomain()
	if err != nil {
		return err
	}

	for _, quantity := range quantities {
		newQuantity := productentity.NewQuantity(productentity.QuantityCommonAttributes{
			Quantity:   quantity,
			CategoryID: *categoryID,
		})

		quantityModel := &model.Quantity{}
		quantityModel.FromDomain(newQuantity)
		if err := s.rq.CreateQuantity(ctx, quantityModel); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetQuantitiesByCategoryId(ctx context.Context, categoryId string) ([]*quantitydto.QuantityDTO, error) {
	if quantities, err := s.rq.GetQuantitiesByCategoryId(ctx, categoryId); err != nil {
		return nil, err
	} else {
		quantityDtos := []*quantitydto.QuantityDTO{}
		for _, quantity := range quantities {
			qDto := &quantitydto.QuantityDTO{}
			qDto.FromDomain(quantity.ToDomain())
			quantityDtos = append(quantityDtos, qDto)
		}
		return quantityDtos, nil
	}
}
