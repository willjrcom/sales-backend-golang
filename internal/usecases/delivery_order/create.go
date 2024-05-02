package deliveryorderusecases

import (
	"context"

	deliveryorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery"
)

func (s *Service) CreateDeliveryOrder(ctx context.Context, dto *deliveryorderdto.CreateDeliveryOrderInput) (*deliveryorderdto.DeliveryIDAndOrderIDOutput, error) {
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

	delivery.AddressID = client.Address.ID
	delivery.DeliveryTax = &client.Address.DeliveryTax

	if err = s.rdo.CreateDeliveryOrder(ctx, delivery); err != nil {
		return nil, err
	}

	return deliveryorderdto.NewOutput(delivery.ID, orderID), nil
}
