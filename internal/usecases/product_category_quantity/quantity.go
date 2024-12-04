package productcategoryquantityusecases

import (
	"context"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategoryquantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category_quantity"
)

type Service struct {
	rq productentity.QuantityRepository
	rc productentity.CategoryRepository
}

func NewService(rq productentity.QuantityRepository) *Service {
	return &Service{rq: rq}
}

func (s *Service) AddDependencies(rc productentity.CategoryRepository) {
	s.rc = rc
}

func (s *Service) CreateQuantity(ctx context.Context, dto *productcategoryquantitydto.CreateQuantityInput) (uuid.UUID, error) {
	quantity, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	category, err := s.rc.GetCategoryById(ctx, quantity.CategoryID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if err = productentity.ValidateDuplicateQuantities(quantity.Quantity, category.Quantities); err != nil {
		return uuid.Nil, err
	}

	if err = s.rq.CreateQuantity(ctx, quantity); err != nil {
		return uuid.Nil, err
	}

	return quantity.ID, nil
}

func (s *Service) UpdateQuantity(ctx context.Context, dtoId *entitydto.IdRequest, dto *productcategoryquantitydto.UpdateQuantityInput) error {
	quantity, err := s.rq.GetQuantityById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(quantity); err != nil {
		return err
	}

	category, err := s.rc.GetCategoryById(ctx, quantity.CategoryID.String())

	if err != nil {
		return err
	}

	if err = productentity.ValidateUpdateQuantity(quantity, category.Quantities); err != nil {
		return err
	}

	if err = s.rq.UpdateQuantity(ctx, quantity); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteQuantity(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.rq.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rq.DeleteQuantity(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQuantityById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.Quantity, error) {
	if quantity, err := s.rq.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return quantity, nil
	}
}

func (s *Service) AddQuantitiesByValues(ctx context.Context, dto *productcategoryquantitydto.RegisterQuantities) error {
	quantities, categoryID, err := dto.ToModel()
	if err != nil {
		return err
	}

	for _, quantity := range quantities {
		newQuantity := productentity.NewQuantity(productentity.QuantityCommonAttributes{
			Quantity:   quantity,
			CategoryID: *categoryID,
		})

		if err := s.rq.CreateQuantity(ctx, newQuantity); err != nil {
			return err
		}
	}

	return nil
}
