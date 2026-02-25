package printmanagerusecases

import (
	"context"

	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
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
	companyRepository   model.CompanyRepository
	rabbitmq            *rabbitmq.RabbitMQ
}

// NewService creates a new print service using the given order and report usecase services.
func NewService() *Service {
	return &Service{}
}

func (s *Service) AddDependencies(orderService *orderusecases.OrderService, orderRepository model.OrderRepository, shiftService *shiftusecases.Service, groupItemRepository model.GroupItemRepository, companyRepository model.CompanyRepository, rabbitmq *rabbitmq.RabbitMQ) {
	s.orderService = orderService
	s.orderRepository = orderRepository
	s.shiftService = shiftService
	s.groupItemRepository = groupItemRepository
	s.companyRepository = companyRepository
	s.rabbitmq = rabbitmq
}

func (s *Service) getCompany(ctx context.Context) (*companydto.CompanyDTO, error) {
	companyModel, err := s.companyRepository.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	company := companyModel.ToDomain()
	dto := &companydto.CompanyDTO{}
	dto.FromDomain(company)
	return dto, nil
}
