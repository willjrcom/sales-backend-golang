package quantityusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
)

var (
	ErrQuantityIsUsed = errors.New("quantity is used in products")
)

type Service struct {
	r productentity.QuantityRepository
}

func NewService(c productentity.QuantityRepository) *Service {
	return &Service{r: c}
}

func (s *Service) RegisterQuantity(ctx context.Context, dto *productdto.RegisterQuantityInput) (uuid.UUID, error) {
	quantity, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.RegisterQuantity(ctx, quantity)

	if err != nil {
		return uuid.Nil, err
	}

	return quantity.ID, nil
}

func (s *Service) UpdateQuantity(ctx context.Context, dtoId *entitydto.IdRequest, dto *productdto.UpdateQuantityInput) error {
	quantity, err := s.r.GetQuantityById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(quantity); err != nil {
		return err
	}

	if err = s.r.UpdateQuantity(ctx, quantity); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteQuantity(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteQuantity(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQuantityById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.Quantity, error) {
	if quantity, err := s.r.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return quantity, nil
	}
}
