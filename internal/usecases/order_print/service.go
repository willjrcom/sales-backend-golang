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
	orderService    *orderusecases.Service
	orderRepository model.OrderRepository
	shiftService    *shiftusecases.Service
}

// NewService creates a new print service using the given order and report usecase services.
func NewService() *Service {
	return &Service{}
}

func (s *Service) AddDependencies(orderService *orderusecases.Service, orderRepository model.OrderRepository, shiftService *shiftusecases.Service) {
	s.orderService = orderService
	s.orderRepository = orderRepository
	s.shiftService = shiftService
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
