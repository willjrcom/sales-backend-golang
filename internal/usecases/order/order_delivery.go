package orderusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrOrderLaunched  = errors.New("order already launched")
	ErrOrderDelivered = errors.New("order already delivered")
)

type IDeliveryService interface {
	ISetupDeliveryService
	ICreateDeliveryService
	IGetDeliveryService
	IUpdateDeliveryService
}
type ISetupDeliveryService interface {
	AddDependencies(ra model.AddressRepository, rc model.ClientRepository, ro model.OrderRepository, so *OrderService, rdd model.DeliveryDriverRepository)
}

type ICreateDeliveryService interface {
	CreateOrderDelivery(ctx context.Context, dto *orderdeliverydto.DeliveryOrderCreateDTO) (*orderdeliverydto.OrderDeliveryIDDTO, error)
}

type IGetDeliveryService interface {
	GetDeliveryById(ctx context.Context, dto *entitydto.IDRequest) (*orderdeliverydto.OrderDeliveryDTO, error)
	GetAllDeliveries(ctx context.Context) ([]orderdeliverydto.OrderDeliveryDTO, error)
}

type IUpdateDeliveryService interface {
	PendOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	ReadyOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	ShipOrderDelivery(ctx context.Context, dtoDriver *orderdeliverydto.DeliveryOrderUpdateShipDTO) (err error)
	CancelOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	OrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IDRequest) (err error)
	UpdateDeliveryDriver(ctx context.Context, dto *entitydto.IDRequest, orderDelivery *orderdeliverydto.DeliveryOrderDriverUpdateDTO) (err error)
	UpdateDeliveryChange(ctx context.Context, dtoId *entitydto.IDRequest, dtoDelivery *orderdeliverydto.OrderChangeCreateDTO) (err error)
}

type OrderDeliveryService struct {
	rdo model.OrderDeliveryRepository
	ra  model.AddressRepository
	rc  model.ClientRepository
	ro  model.OrderRepository
	so  *OrderService
	rdd model.DeliveryDriverRepository
}

func NewDeliveryService(rdo model.OrderDeliveryRepository) IDeliveryService {
	return &OrderDeliveryService{rdo: rdo}
}

func (s *OrderDeliveryService) AddDependencies(ra model.AddressRepository, rc model.ClientRepository, ro model.OrderRepository, os *OrderService, rdd model.DeliveryDriverRepository) {
	s.ra = ra
	s.rc = rc
	s.ro = ro
	s.so = os
	s.rdd = rdd
}

func (s *OrderDeliveryService) CreateOrderDelivery(ctx context.Context, dto *orderdeliverydto.DeliveryOrderCreateDTO) (*orderdeliverydto.OrderDeliveryIDDTO, error) {
	delivery, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	orderID, err := s.so.CreateDefaultOrder(ctx)

	if err != nil {
		return nil, err
	}

	delivery.OrderID = orderID

	// Validate client
	clientModel, err := s.rc.GetClientById(ctx, delivery.ClientID.String())
	if err != nil {
		return nil, err
	}

	client := clientModel.ToDomain()

	delivery.ClientID = client.ID
	delivery.AddressID = client.Address.ID
	delivery.DeliveryTax = &client.Address.DeliveryTax

	deliveryModel := &model.OrderDelivery{}
	deliveryModel.FromDomain(delivery)
	if err = s.rdo.CreateOrderDelivery(ctx, deliveryModel); err != nil {
		return nil, err
	}

	// Update delivery tax
	if err := s.so.UpdateOrderTotal(ctx, delivery.OrderID.String()); err != nil {
		return nil, err
	}

	return orderdeliverydto.FromDomain(delivery.ID, orderID), nil
}

func (s *OrderDeliveryService) GetDeliveryById(ctx context.Context, dto *entitydto.IDRequest) (*orderdeliverydto.OrderDeliveryDTO, error) {
	if deliveryModel, err := s.rdo.GetDeliveryById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		delivery := deliveryModel.ToDomain()

		deliveryDTO := &orderdeliverydto.OrderDeliveryDTO{}
		deliveryDTO.FromDomain(delivery)
		return deliveryDTO, nil
	}
}

func (s *OrderDeliveryService) GetAllDeliveries(ctx context.Context) ([]orderdeliverydto.OrderDeliveryDTO, error) {
	if deliveryModels, err := s.rdo.GetAllDeliveries(ctx); err != nil {
		return nil, err
	} else {
		deliveries := []orderdeliverydto.OrderDeliveryDTO{}
		for _, deliveryModel := range deliveryModels {
			delivery := deliveryModel.ToDomain()
			deliveryDTO := &orderdeliverydto.OrderDeliveryDTO{}
			deliveryDTO.FromDomain(delivery)
			deliveries = append(deliveries, *deliveryDTO)
		}
		return deliveries, nil
	}
}

func (s *OrderDeliveryService) PendOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
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
func (s *OrderDeliveryService) ReadyOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
	orderDeliveryModel, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderDelivery := orderDeliveryModel.ToDomain()
	if err := orderDelivery.Ready(); err != nil {
		return err
	}

	orderDeliveryModel.FromDomain(orderDelivery)
	if err = s.rdo.UpdateOrderDelivery(ctx, orderDeliveryModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderDeliveryService) ShipOrderDelivery(ctx context.Context, dtoShip *orderdeliverydto.DeliveryOrderUpdateShipDTO) (err error) {
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

func (s *OrderDeliveryService) CancelOrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
	orderDeliveryModel, err := s.rdo.GetDeliveryById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	orderDelivery := orderDeliveryModel.ToDomain()
	if err := orderDelivery.Cancel(); err != nil {
		return err
	}

	orderDeliveryModel.FromDomain(orderDelivery)
	if err = s.rdo.UpdateOrderDelivery(ctx, orderDeliveryModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderDeliveryService) OrderDelivery(ctx context.Context, dtoID *entitydto.IDRequest) (err error) {
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

func (s *OrderDeliveryService) UpdateDeliveryAddress(ctx context.Context, dtoID *entitydto.IDRequest) error {
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

func (s *OrderDeliveryService) UpdateDeliveryDriver(ctx context.Context, dtoID *entitydto.IDRequest, dto *orderdeliverydto.DeliveryOrderDriverUpdateDTO) error {
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

func (s *OrderDeliveryService) UpdateDeliveryChange(ctx context.Context, dto *entitydto.IDRequest, dtoPayment *orderdeliverydto.OrderChangeCreateDTO) error {
	change, method, err := dtoPayment.ToDomain()
	if err != nil {
		return err
	}

	orderDeliveryModel, err := s.rdo.GetDeliveryById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	orderDelivery := orderDeliveryModel.ToDomain()

	orderDelivery.AddChange(change, *method)

	orderDeliveryModel.FromDomain(orderDelivery)
	if err := s.rdo.UpdateOrderDelivery(ctx, orderDeliveryModel); err != nil {
		return err
	}

	return nil
}
