package deliveryorderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) UpdateDeliveryAddress(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.UpdateDeliveryOrder) error {
	_, err := s.ra.GetAddressById(ctx, dto.AddressID.String())

	if err != nil {
		return err
	}

	return nil
}
