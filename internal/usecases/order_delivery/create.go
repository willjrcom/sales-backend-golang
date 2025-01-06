package orderdeliveryusecases

import (
	"context"

	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *Service) CreateOrderDelivery(ctx context.Context, dto *orderdeliverydto.DeliveryOrderCreateDTO) (*orderdeliverydto.OrderDeliveryIDDTO, error) {
	delivery, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	delivery.OrderID = orderID

	// Validate client
	clientModel, err := s.rc.GetClientById(ctx, delivery.ClientID.String())
	if err != nil {
		return nil, err
	}

	client := clientModel.ToDomain()

	delivery.ClientID = client.ID
	delivery.AddressID = client.Address.ID
	delivery.DeliveryTax = &client.Address.DeliveryTax

	deliveryModel := &model.OrderDelivery{}
	deliveryModel.FromDomain(delivery)
	if err = s.rdo.CreateOrderDelivery(ctx, deliveryModel); err != nil {
		return nil, err
	}

	return orderdeliverydto.FromDomain(delivery.ID, orderID), nil
}
