package orderusecases

import (
	"context"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

type Service struct {
	r orderentity.Repository
}

func NewService(repository orderentity.Repository) *Service {
	return &Service{r: repository}
}

func (s *Service) CreateOrder(ctx context.Context, dto *orderdto.CreateOrderInput) (uuid.UUID, error) {
	order, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.r.CreateOrder(ctx, order); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}

func (s *Service) LaunchOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.r.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order.Status = orderentity.OrderStatusPending

	if err := s.r.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) ArchiveOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.r.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order.Status = orderentity.OrderStatusArchived

	if err := s.r.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelOrder(ctx context.Context, dto *entitydto.IdRequest) error {
	order, err := s.r.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order.Status = orderentity.OrderStatusCanceled

	if err := s.r.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateOrderPayment(ctx context.Context, dto *entitydto.IdRequest, dtoPayment *orderdto.UpdatePaymentMethod) error {
	order, err := s.r.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	dtoPayment.UpdateModel(order)

	if err := s.r.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateOrderObservation(ctx context.Context, dtoId *entitydto.IdRequest, dto *orderdto.UpdateObservationOrder) error {
	order, err := s.r.GetOrderById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	dto.UpdateModel(order)

	if err := s.r.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryOrder(ctx context.Context) error {
	return nil
}

func (s *Service) UpdateTableOrder(ctx context.Context) error {
	return nil
}

func (s *Service) UpdateOrderStatus(ctx context.Context) error {
	return nil
}

func (s *Service) GetOrderById(ctx context.Context, dto *entitydto.IdRequest) (*orderentity.Order, error) {
	if order, err := s.r.GetOrderById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return order, nil
	}
}

func (s *Service) GetAllOrder(ctx context.Context) ([]orderentity.Order, error) {
	if orders, err := s.r.GetAllOrders(ctx); err != nil {
		return nil, err
	} else {
		return orders, nil
	}
}
