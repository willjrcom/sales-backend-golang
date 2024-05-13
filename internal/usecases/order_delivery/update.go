package orderdeliveryusecases

import (
	"context"
	"errors"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
)

var (
	ErrOrderLaunched  = errors.New("order already launched")
	ErrOrderDelivered = errors.New("order already delivered")
)

func (s *Service) PendOrderDelivery(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	orderDelivery, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := orderDelivery.Pend(); err != nil {
		return err
	}

	if err = s.rdo.UpdateOrderDelivery(ctx, orderDelivery); err != nil {
		return err
	}

	return nil
}

func (s *Service) ShipOrderDelivery(ctx context.Context, dtoID *entitydto.IdRequest, dtoDriver *orderdeliverydto.UpdateDriverOrder) (err error) {
	orderDelivery, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err = dtoDriver.UpdateModel(orderDelivery); err != nil {
		return err
	}

	if _, err = s.re.GetEmployeeById(ctx, orderDelivery.DriverID.String()); err != nil {
		return err
	}

	if err := orderDelivery.Ship(orderDelivery.DriverID); err != nil {
		return err
	}

	if err = s.rdo.UpdateOrderDelivery(ctx, orderDelivery); err != nil {
		return err
	}

	return nil
}

func (s *Service) OrderDelivery(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	orderDelivery, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := orderDelivery.Delivery(); err != nil {
		return err
	}

	if err = s.rdo.UpdateOrderDelivery(ctx, orderDelivery); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IdRequest) error {
	orderDelivery, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if orderDelivery.DeliveredAt != nil {
		return ErrOrderDelivered
	}

	if orderDelivery.ShippedAt != nil {
		return ErrOrderLaunched
	}

	orderDelivery.AddressID = orderDelivery.Client.Address.ID

	if err := s.rdo.UpdateOrderDelivery(ctx, orderDelivery); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryDriver(ctx context.Context, dtoID *entitydto.IdRequest, dto *orderdeliverydto.UpdateDriverOrder) error {
	orderDelivery, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if orderDelivery.DeliveredAt != nil {
		return ErrOrderDelivered
	}

	if _, err = s.re.GetEmployeeById(ctx, dto.DriverID.String()); err != nil {
		return err
	}

	if err := dto.UpdateModel(orderDelivery); err != nil {
		return err
	}

	if err := s.rdo.UpdateOrderDelivery(ctx, orderDelivery); err != nil {
		return err
	}

	return nil
}
