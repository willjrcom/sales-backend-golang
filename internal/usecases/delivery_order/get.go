package deliveryorderusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetDeliveryById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.DeliveryOrder, error) {
	if delivery, err := s.rdo.GetDeliveryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return delivery, nil
	}
}

func (s *Service) GetAllDeliveries(ctx context.Context) ([]orderentity.DeliveryOrder, error) {
	if deliveries, err := s.rdo.GetAllDeliveries(ctx); err != nil {
		return nil, err
	} else {
		return deliveries, nil
	}
}
