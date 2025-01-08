package orderdeliveryusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
)

var (
	ErrOrderLaunched  = errors.New("order already launched")
	ErrOrderDelivered = errors.New("order already delivered")
)

func (s *Service) PendOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
	orderDeliveryModel, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderDelivery := orderDeliveryModel.ToDomain()
	if err := orderDelivery.Pend(); err != nil {
		return err
	}

	orderDeliveryModel.FromDomain(orderDelivery)
	if err = s.rdo.UpdateOrderDelivery(ctx, orderDeliveryModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) ShipOrderDelivery(ctx context.Context, dtoShip *orderdeliverydto.DeliveryOrderUpdateShipDTO) (err error) {
	if len(dtoShip.DeliveryIDs) == 0 {
		return errors.New("delivery ids is required")
	}

	orderDeliveryModel, err := s.rdo.GetDeliveriesByIds(ctx, dtoShip.DeliveryIDs)

	if err != nil {
		return err
	}

	orderDeliveries := []orderentity.OrderDelivery{}
	for _, orderDeliveryModel := range orderDeliveryModel {
		orderDeliveries = append(orderDeliveries, *orderDeliveryModel.ToDomain())
	}

	if err = dtoShip.UpdateDomain(orderDeliveries); err != nil {
		return err
	}

	if _, err = s.rdd.GetDeliveryDriverById(ctx, dtoShip.DriverID.String()); err != nil {
		return err
	}

	for i := range orderDeliveries {
		if err := orderDeliveries[i].Ship(&dtoShip.DriverID); err != nil {
			return err
		}

		orderDeliveryModel[i].FromDomain(&orderDeliveries[i])
		if err := s.rdo.UpdateOrderDelivery(ctx, &orderDeliveryModel[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) OrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
	orderDeliveryModel, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderDelivery := orderDeliveryModel.ToDomain()
	if err := orderDelivery.Delivery(); err != nil {
		return err
	}

	orderDeliveryModel.FromDomain(orderDelivery)
	if err = s.rdo.UpdateOrderDelivery(ctx, orderDeliveryModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IDRequest) error {
	orderDeliveryModel, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if orderDeliveryModel.DeliveredAt != nil {
		return ErrOrderDelivered
	}

	if orderDeliveryModel.ShippedAt != nil {
		return ErrOrderLaunched
	}

	orderDeliveryModel.AddressID = orderDeliveryModel.Client.Address.ID

	if err := s.rdo.UpdateOrderDelivery(ctx, orderDeliveryModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryDriver(ctx context.Context, dtoID *entitydto.IDRequest, dto *orderdeliverydto.DeliveryOrderDriverUpdateDTO) error {
	orderDeliveryModel, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderDelivery := orderDeliveryModel.ToDomain()
	if orderDelivery.DeliveredAt != nil {
		return ErrOrderDelivered
	}

	if _, err = s.rdd.GetDeliveryDriverById(ctx, dto.DriverID.String()); err != nil {
		return err
	}

	if err := dto.UpdateDomain(orderDelivery); err != nil {
		return err
	}

	orderDeliveryModel.FromDomain(orderDelivery)
	if err := s.rdo.UpdateOrderDelivery(ctx, orderDeliveryModel); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, orderDelivery.OrderID.String()); err != nil {
		return err
	}

	return nil
}
