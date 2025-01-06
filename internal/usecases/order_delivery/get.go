package orderdeliveryusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) GetDeliveryById(ctx context.Context, dto *entitydto.IDRequest) (*orderentity.OrderDelivery, error) {
	if deliveryModel, err := s.rdo.GetDeliveryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return deliveryModel.ToDomain(), nil
	}
}

func (s *Service) GetAllDeliveries(ctx context.Context) ([]orderentity.OrderDelivery, error) {
	if deliveryModels, err := s.rdo.GetAllDeliveries(ctx); err != nil {
		return nil, err
	} else {
		deliveries := []orderentity.OrderDelivery{}
		for _, deliveryModel := range deliveryModels {
			deliveries = append(deliveries, *deliveryModel.ToDomain())
		}
		return deliveries, nil
	}
}

func (s *Service) GetAllOrderDeliveryStatus(ctx context.Context) (deliveries []orderentity.StatusOrderDelivery) {
	return orderentity.GetAllDeliveryStatus()
}

func (s Service) GetOrderDeliveryByClientId(ctx context.Context, dto *entitydto.IDRequest) ([]orderentity.OrderDelivery, error) {
	return nil, nil
}

func (s *Service) GetOrderDeliveryByDriverId(ctx context.Context, dto *entitydto.IDRequest) ([]orderentity.OrderDelivery, error) {
	return nil, nil
}

func (s *Service) GetOrderDeliveryByStatus(ctx context.Context) (deliveries []orderentity.OrderDelivery, err error) {
	return nil, nil
}
