package deliveryorderusecases

import (
	"context"
	"errors"

	deliveryorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

var (
	ErrOrderLaunched  = errors.New("order already launched")
	ErrOrderDelivered = errors.New("order already delivered")
)

func (s *Service) PendDeliveryOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	deliveryOrder, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := deliveryOrder.Pend(); err != nil {
		return err
	}

	if err = s.rdo.UpdateDeliveryOrder(ctx, deliveryOrder); err != nil {
		return err
	}

	return nil
}

func (s *Service) ShipDeliveryOrder(ctx context.Context, dtoID *entitydto.IdRequest, dtoDriver *deliveryorderdto.UpdateDriverOrder) (err error) {
	deliveryOrder, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err = dtoDriver.UpdateModel(deliveryOrder); err != nil {
		return err
	}

	if _, err = s.re.GetEmployeeById(ctx, deliveryOrder.DriverID.String()); err != nil {
		return err
	}

	if err := deliveryOrder.Ship(*deliveryOrder.DriverID); err != nil {
		return err
	}

	if err = s.rdo.UpdateDeliveryOrder(ctx, deliveryOrder); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeliveryOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	deliveryOrder, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := deliveryOrder.Delivery(); err != nil {
		return err
	}

	if err = s.rdo.UpdateDeliveryOrder(ctx, deliveryOrder); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IdRequest) error {
	deliveryOrder, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if deliveryOrder.DeliveredAt != nil {
		return ErrOrderDelivered
	}

	if deliveryOrder.ShippedAt != nil {
		return ErrOrderLaunched
	}

	deliveryOrder.AddressID = deliveryOrder.Client.Address.ID

	if err := s.rdo.UpdateDeliveryOrder(ctx, deliveryOrder); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryDriver(ctx context.Context, dtoID *entitydto.IdRequest, dto *deliveryorderdto.UpdateDriverOrder) error {
	deliveryOrder, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if deliveryOrder.DeliveredAt != nil {
		return ErrOrderDelivered
	}

	if _, err = s.re.GetEmployeeById(ctx, dto.DriverID.String()); err != nil {
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
