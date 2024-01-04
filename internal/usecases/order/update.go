package orderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

func (s *Service) PendingOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.PendingOrder(); err != nil {
		return err
	}

	if err := s.ro.PendingOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) FinishOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.FinishOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelOrder(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.CancelOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) ArchiveOrder(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.ArchiveOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UnarchiveOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.UnarchiveOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdatePaymentMethod(ctx context.Context, dto *entitydto.IdRequest, dtoPayment *orderdto.UpdatePaymentMethod) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	dtoPayment.UpdateModel(order)

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateOrderObservation(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.UpdateObservationOrder) error {
	order, err := s.ro.GetOrderById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	dto.UpdateModel(order)

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}
