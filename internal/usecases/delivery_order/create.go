package deliveryorderusecases

import (
	"context"

	"github.com/google/uuid"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) CreateDeliveryOrder(ctx context.Context, dto *orderdto.CreateDeliveryOrderInput) (uuid.UUID, error) {
	delivery, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	// Validate address
	if _, err := s.ra.GetAddressById(ctx, delivery.AddressID.String()); err != nil {
		return uuid.Nil, err
	}

	// Validate client
	if _, err := s.rc.GetClientById(ctx, delivery.ClientID.String()); err != nil {
		return uuid.Nil, err
	}

	// Validate order
	if _, err := s.ro.GetOrderById(ctx, delivery.OrderID.String()); err != nil {
		return uuid.Nil, err
	}

	if err = s.rdo.CreateDeliveryOrder(ctx, delivery); err != nil {
		return uuid.Nil, err
	}

	return delivery.ID, nil
}
