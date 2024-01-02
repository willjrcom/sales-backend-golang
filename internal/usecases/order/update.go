package orderusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

var (
	ErrOrderWithoutItems          = errors.New("order must have at least one item")
	ErrOrderNotCanceledOrFinished = errors.New("order must be canceled or finished")
	ErrOrderAlreadyCanceled       = errors.New("order already canceled")
	ErrOrderAlreadyArchived       = errors.New("order already archived")
)

func (s *Service) LaunchOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if len(order.Groups) == 0 {
		return ErrOrderWithoutItems
	}

	order.LaunchOrder()

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) ArchiveOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if order.Status != orderentity.OrderStatusCanceled && order.Status != orderentity.OrderStatusFinished {
		return ErrOrderNotCanceledOrFinished
	}

	order.ArchiveOrder()

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if order.Status == orderentity.OrderStatusCanceled {
		return ErrOrderAlreadyCanceled
	}

	if order.Status == orderentity.OrderStatusArchived {
		return ErrOrderAlreadyArchived
	}
	order.CancelOrder()

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	// cancelar todos itens
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
