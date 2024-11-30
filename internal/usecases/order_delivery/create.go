package orderdeliveryusecases

import (
	"context"

	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
)

func (s *Service) CreateOrderDelivery(ctx context.Context, dto *orderdeliverydto.CreateOrderDeliveryInput) (*orderdeliverydto.DeliveryIDAndOrderIDOutput, error) {
	delivery, err := dto.ToModel()

	if err != nil {
		return nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	delivery.OrderID = orderID

	// Validate client
	client, err := s.rc.GetClientById(ctx, delivery.ClientID.String())
	if err != nil {
		return nil, err
	}

	delivery.ClientID = client.ID
	delivery.AddressID = client.Address.ID
	delivery.DeliveryTax = &client.Address.DeliveryTax

	if err = s.rdo.CreateOrderDelivery(ctx, delivery); err != nil {
		return nil, err
	}

	return orderdeliverydto.NewOutput(delivery.ID, orderID), nil
}
