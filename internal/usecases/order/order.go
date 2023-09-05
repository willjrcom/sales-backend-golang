package orderusecases

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
)

type Service struct {
	Repository orderentity.Repository
}

func NewService(repository orderentity.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) CreateOrder(dto *orderdto.CreateOrderInput) (uuid.UUID, error) {
	order, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.Repository.CreateOrder(order); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}

func (s *Service) LaunchOrder(dto *entitydto.IdRequest) error {
	order, err := s.Repository.GetOrderById(dto.Id.String())

	if err != nil {
		return err
	}

	order.Status = orderentity.OrderStatusPending

	if err := s.Repository.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}

func (s *Service) ArchiveOrder(dto *entitydto.IdRequest) error {
	order, err := s.Repository.GetOrderById(dto.Id.String())

	if err != nil {
		return err
	}

	order.Status = orderentity.OrderStatusArchived

	if err := s.Repository.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelOrder(dto *entitydto.IdRequest) error {
	order, err := s.Repository.GetOrderById(dto.Id.String())

	if err != nil {
		return err
	}

	order.Status = orderentity.OrderStatusCanceled

	if err := s.Repository.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateOrderPayment(dto *entitydto.IdRequest, dtoPayment *orderdto.UpdatePaymentMethod) error {
	order, err := s.Repository.GetOrderById(dto.Id.String())

	if err != nil {
		return err
	}

	dtoPayment.UpdateModel(order)

	if err := s.Repository.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateOrderObservation(dtoId *entitydto.IdRequest, dto *orderdto.UpdateObservationOrder) error {
	order, err := s.Repository.GetOrderById(dtoId.Id.String())

	if err != nil {
		return err
	}

	dto.UpdateModel(order)

	if err := s.Repository.UpdateOrder(order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDeliveryOrder() error {
	return nil
}

func (s *Service) UpdateTableOrder() error {
	return nil
}

func (s *Service) UpdateOrderStatus() error {
	return nil
}

func (s *Service) GetOrderById(dto *entitydto.IdRequest) (*orderentity.Order, error) {
	if order, err := s.Repository.GetOrderById(dto.Id.String()); err != nil {
		return nil, err
	} else {
		return order, nil
	}
}

func (s *Service) GetAllOrder(dto *filterdto.Filter) ([]orderentity.Order, error) {
	if orders, err := s.Repository.GetAllOrder(dto.Key, dto.Value); err != nil {
		return nil, err
	} else {
		return orders, nil
	}
}
