package deliveryorderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) UpdateDeliveryAddress(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.UpdateDeliveryOrder) error {
	deliveryOrder, err := s.rdo.GetDeliveryById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	address, err := s.ra.GetAddressById(ctx, dto.AddressID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(deliveryOrder, address); err != nil {
		return err
	}

	if err := s.rdo.UpdateDeliveryOrder(ctx, deliveryOrder); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDriver(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.UpdateDriverOrder) error {
	deliveryOrder, err := s.rdo.GetDeliveryById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	_, err = s.re.GetEmployeeById(ctx, dto.DriverID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(deliveryOrder); err != nil {
		return err
	}

	if err := s.rdo.UpdateDeliveryOrder(ctx, deliveryOrder); err != nil {
		return err
	}

	return nil
}
