package deliveryorderusecases

import (
	"context"

	"github.com/google/uuid"
	deliveryorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery"
)

func (s *Service) CreateDeliveryOrder(ctx context.Context, dto *deliveryorderdto.CreateDeliveryOrderInput) (uuid.UUID, error) {
	delivery, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	orderID, err := s.os.CreateDefaultOrder(ctx)

	if err != nil {
		return uuid.Nil, err
	}

	delivery.OrderID = orderID

	// Validate client
	client, err := s.rc.GetClientById(ctx, delivery.ClientID.String())
	if err != nil {
		return uuid.Nil, err
	}

	delivery.AddressID = client.Address.ID

	if err = s.rdo.CreateDeliveryOrder(ctx, delivery); err != nil {
		return uuid.Nil, err
	}

	return delivery.ID, nil
}
