package orderusecases

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
	orderprocessusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_process"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

type Service struct {
	ro  orderentity.OrderRepository
	rs  shiftentity.ShiftRepository
	rgi *groupitemusecases.Service
	rp  *orderprocessusecases.Service
	rpr productentity.ProcessRuleRepository
	rq  *orderqueueusecases.Service
}

func NewService(ro orderentity.OrderRepository) *Service {
	return &Service{ro: ro}
}

func (s *Service) AddDependencies(rs shiftentity.ShiftRepository, rgi *groupitemusecases.Service, rp *orderprocessusecases.Service, rpr productentity.ProcessRuleRepository, rq *orderqueueusecases.Service) {
	s.rs = rs
	s.rgi = rgi
	s.rp = rp
	s.rpr = rpr
	s.rq = rq
}
