package orderpickupusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
)

func (s *Service) PendingOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	orderPickup, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := orderPickup.Pend(); err != nil {
		return err
	}

	if err = s.rp.UpdateOrderPickup(ctx, orderPickup); err != nil {
		return err
	}

	return nil
}

func (s *Service) ReadyOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	orderPickup, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := orderPickup.Ready(); err != nil {
		return err
	}

	if err = s.rp.UpdateOrderPickup(ctx, orderPickup); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateName(ctx context.Context, dtoID *entitydto.IdRequest, dtoPickup *orderpickupdto.UpdateOrderPickupInput) (err error) {
	orderPickup, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := orderPickup.UpdateName(dtoPickup.Name); err != nil {
		return err
	}

	if err = s.rp.UpdateOrderPickup(ctx, orderPickup); err != nil {
		return err
	}

	return nil
}
