package orderusecases

import (
	"github.com/google/uuid"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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
	// save order
	return order.ID, nil
}

func (s *Service) AddItemOrder(idOrder, idProduct string, dto *orderdto.AddItemOrderInput) (uuid.UUID, error) {
	// find order and product by id
	items := itementity.Items{}
	product := &productentity.Product{}
	item := dto.ToModel(product)

	item.ItemsID = items.ID
	items.Products = append(items.Products, *item)
	return uuid.New(), nil
}

func (s *Service) LaunchOrder(dto *entitydto.IdRequest) error {
	order, err := s.Repository.GetOrder(dto.Id.String())

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
	order, err := s.Repository.GetOrder(dto.Id.String())

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
	order, err := s.Repository.GetOrder(dto.Id.String())

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
	order, err := s.Repository.GetOrder(dto.Id.String())

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
	order, err := s.Repository.GetOrder(dtoId.Id.String())

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

func (s *Service) GetOrder(dto *entitydto.IdRequest) (*orderentity.Order, error) {
	if order, err := s.Repository.GetOrder(dto.Id.String()); err != nil {
		return nil, err
	} else {
		return order, nil
	}
}

func (s *Service) GetAllOrder() ([]orderentity.Order, error) {
	if orders, err := s.Repository.GetAllOrder(); err != nil {
		return nil, err
	} else {
		return orders, nil
	}
}
