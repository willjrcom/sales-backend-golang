package orderprintusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/pos"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
	shiftusecases "github.com/willjrcom/sales-backend-go/internal/usecases/shift"
)

// Service provides print operations for orders.
// Service provides print operations for orders and daily reports.
type Service struct {
	orderService        *orderusecases.OrderService
	shiftService        *shiftusecases.Service
	orderRepository     model.OrderRepository
	groupItemRepository model.GroupItemRepository
}

// NewService creates a new print service using the given order and report usecase services.
func NewService() *Service {
	return &Service{}
}

func (s *Service) AddDependencies(orderService *orderusecases.OrderService, orderRepository model.OrderRepository, shiftService *shiftusecases.Service, groupItemRepository model.GroupItemRepository) {
	s.orderService = orderService
	s.orderRepository = orderRepository
	s.shiftService = shiftService
	s.groupItemRepository = groupItemRepository
}

// PrintOrder retrieves the order by ID and returns its printable representation.
func (s *Service) PrintOrder(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	model, err := s.orderRepository.GetOrderById(ctx, req.ID.String())
	if err != nil {
		return nil, err
	}

	order := model.ToDomain()
	data, err := pos.FormatOrder(order)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// PrintDailyReport retrieves daily sales summary for a specific day.
func (s *Service) PrintDailyReport(ctx context.Context, req *entitydto.IDRequest) (interface{}, error) {
	shift, err := s.shiftService.GetShiftByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return shift, nil
}

// PrintGroupItemKitchen retrieves the order by ID and returns its kitchen-printable bytes
// showing only items and complements, without prices or totals.
func (s *Service) PrintGroupItemKitchen(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	// fetch full order model
	modelGroupItem, err := s.groupItemRepository.GetGroupByID(ctx, req.ID.String(), true)
	if err != nil {
		return nil, err
	}

	// convert to domain
	groupItem := modelGroupItem.ToDomain()
	data, err := pos.FormatGroupItemKitchen(groupItem)
	if err != nil {
		return nil, err
	}
	return data, nil
}
