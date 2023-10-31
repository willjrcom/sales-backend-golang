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

func (s *Service) GetAllDeliveryOrderStatus(ctx context.Context) (deliveries []orderentity.StatusDeliveryOrder) {
	return orderentity.GetAllDeliveryStatus()
}

func (s Service) GetDeliveryOrderByClientId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.DeliveryOrder, error) {
	return nil, nil
}

func (s *Service) GetDeliveryOrderByDriverId(ctx context.Context, dto *entitydto.IdRequest) ([]orderentity.DeliveryOrder, error) {
	return nil, nil
}

func (s *Service) GetDeliveryOrderByStatus(ctx context.Context) (deliveries []orderentity.DeliveryOrder, err error) {
	return nil, nil
}
