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
}

func NewService() *Service {
	return &Service{}
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

func (s *Service) LaunchOrder(dto *orderdto.LaunchOrderInput) error {
	// find order by id
	order := orderentity.Order{}
	order.Status = orderentity.OrderStatusPending
	// save(order.ID)
	return nil
}

func (s *Service) ArchiveOrder(dto *entitydto.IdRequest) error {
	return nil
}

func (s *Service) CancelOrder(dto *entitydto.IdRequest) error {
	return nil
}

func (s *Service) UpdateDeliveryOrder() error {
	return nil
}

func (s *Service) UpdateOrderPayment() error {
	return nil
}

func (s *Service) UpdateOrderStatus() error {
	return nil
}

func (s *Service) UpdateOrderObservation() error {
	return nil
}

func (s *Service) UpdateTableOrder() error {
	return nil
}

func (s *Service) GetOrder(dto *entitydto.IdRequest) (*orderentity.Order, error) {
	return nil, nil
}

func (s *Service) GetAllOrder() ([]orderentity.Order, error) {
	return nil, nil
}
