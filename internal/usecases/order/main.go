package orderusecases

import (
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

type OrderService struct {
	ro  model.OrderRepository
	rs  model.ShiftRepository
	rp  model.ProductRepository
	rpr model.ProcessRuleRepository
	rdo model.OrderDeliveryRepository
	sgi *GroupItemService
	sop *OrderProcessService
	sq  *orderqueueusecases.Service
	sd  IDeliveryService
	sp  IPickupService
	st  *OrderTableService
	sc  *companyusecases.Service
}

func NewOrderService(ro model.OrderRepository) *OrderService {
	return &OrderService{ro: ro}
}

func (s *OrderService) AddDependencies(
	ro model.OrderRepository,
	rs model.ShiftRepository,
	rp model.ProductRepository,
	rpr model.ProcessRuleRepository,
	rdo model.OrderDeliveryRepository,
	sgi *GroupItemService,
	sop *OrderProcessService,
	sq *orderqueueusecases.Service,
	sd IDeliveryService,
	sp IPickupService,
	st *OrderTableService,
	sc *companyusecases.Service,
) {
	s.ro = ro
	s.rs = rs
	s.rp = rp
	s.rpr = rpr
	s.rdo = rdo
	s.sgi = sgi
	s.sq = sq
	s.sop = sop
	s.sd = sd
	s.sp = sp
	s.st = st
	s.sc = sc
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
