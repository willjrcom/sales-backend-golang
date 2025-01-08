package orderdeliveryusecases

import (
	"context"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
)

func (s *Service) GetDeliveryById(ctx context.Context, dto *entitydto.IDRequest) (*orderdeliverydto.OrderDeliveryDTO, error) {
	if deliveryModel, err := s.rdo.GetDeliveryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		delivery := deliveryModel.ToDomain()

		deliveryDTO := &orderdeliverydto.OrderDeliveryDTO{}
		deliveryDTO.FromDomain(delivery)
		return deliveryDTO, nil
	}
}

func (s *Service) GetAllDeliveries(ctx context.Context) ([]orderdeliverydto.OrderDeliveryDTO, error) {
	if deliveryModels, err := s.rdo.GetAllDeliveries(ctx); err != nil {
		return nil, err
	} else {
		deliveries := []orderdeliverydto.OrderDeliveryDTO{}
		for _, deliveryModel := range deliveryModels {
			delivery := deliveryModel.ToDomain()
			deliveryDTO := &orderdeliverydto.OrderDeliveryDTO{}
			deliveryDTO.FromDomain(delivery)
			deliveries = append(deliveries, *deliveryDTO)
		}
		return deliveries, nil
	}
}

func (s *Service) GetAllOrderDeliveryStatus(ctx context.Context) (deliveries []orderentity.StatusOrderDelivery) {
	return orderentity.GetAllDeliveryStatus()
}

func (s Service) GetOrderDeliveryByClientId(ctx context.Context, dto *entitydto.IDRequest) ([]orderdeliverydto.OrderDeliveryDTO, error) {
	return nil, nil
}

func (s *Service) GetOrderDeliveryByDriverId(ctx context.Context, dto *entitydto.IDRequest) ([]orderdeliverydto.OrderDeliveryDTO, error) {
	return nil, nil
}

func (s *Service) GetOrderDeliveryByStatus(ctx context.Context) (deliveries []orderdeliverydto.OrderDeliveryDTO, err error) {
	return nil, nil
}
