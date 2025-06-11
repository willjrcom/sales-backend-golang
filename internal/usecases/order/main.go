package orderusecases

import (
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

type OrderService struct {
	ro  model.OrderRepository
	rs  model.ShiftRepository
	rp  model.ProductRepository
	rpr model.ProcessRuleRepository
	sgi *GroupItemService
	sop *OrderProcessService
	sq  *orderqueueusecases.Service
}

func NewService(ro model.OrderRepository) *OrderService {
	return &OrderService{ro: ro}
}

func (s *OrderService) AddDependencies(
	ro model.OrderRepository,
	rs model.ShiftRepository,
	rp model.ProductRepository,
	rpr model.ProcessRuleRepository,
	sgi *GroupItemService,
	sop *OrderProcessService,
	sq *orderqueueusecases.Service,
) {
	s.ro = ro
	s.rs = rs
	s.rp = rp
	s.rpr = rpr
	s.sgi = sgi
	s.sq = sq
	s.sop = sop
}

type GroupItemService struct {
	r   model.GroupItemRepository
	ri  model.ItemRepository
	rp  model.ProductRepository
	sop *OrderProcessService
	so  *OrderService
}

func NewGroupItemService(rgi model.GroupItemRepository) *GroupItemService {
	return &GroupItemService{r: rgi}
}

func (s *GroupItemService) AddDependencies(ri model.ItemRepository, rp model.ProductRepository, so *OrderService, sop *OrderProcessService) {
	s.ri = ri
	s.rp = rp
	s.so = so
	s.sop = sop
}

type OrderProcessService struct {
	r   model.OrderProcessRepository
	rpr model.ProcessRuleRepository
	sq  *orderqueueusecases.Service
	sgi *GroupItemService
	rgi model.GroupItemRepository
	ro  model.OrderRepository
	se  *employeeusecases.Service
}

func NewOrderProcessService(c model.OrderProcessRepository) *OrderProcessService {
	return &OrderProcessService{r: c}
}

func (s *OrderProcessService) AddDependencies(sq *orderqueueusecases.Service, rpr model.ProcessRuleRepository, sgi *GroupItemService, ro model.OrderRepository, se *employeeusecases.Service, rgi model.GroupItemRepository) {
	s.rgi = rgi
	s.rpr = rpr
	s.sq = sq
	s.sgi = sgi
	s.ro = ro
	s.se = se
}
