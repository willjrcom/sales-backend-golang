package orderusecases

import (
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

type OrderService struct {
	ro                model.OrderRepository
	rs                model.ShiftRepository
	rp                model.ProductRepository
	rpr               model.ProcessRuleRepository
	rdo               model.OrderDeliveryRepository
	stockRepo         model.StockRepository
	stockMovementRepo model.StockMovementRepository
	sgi               *GroupItemService
	sop               *OrderProcessService
	sq                *orderqueueusecases.Service
	sd                IDeliveryService
	sp                IPickupService
	st                *OrderTableService
	sc                *companyusecases.Service
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
	stockRepo model.StockRepository,
	stockMovementRepo model.StockMovementRepository,
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
	s.stockRepo = stockRepo
	s.stockMovementRepo = stockMovementRepo
	s.sgi = sgi
	s.sq = sq
	s.sop = sop
	s.sd = sd
	s.sp = sp
	s.st = st
	s.sc = sc
}
