package orderusecases

import (
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
	orderprocessusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_process"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

type Service struct {
	ro  model.OrderRepository
	rs  model.ShiftRepository
	rgi *groupitemusecases.Service
	rp  *orderprocessusecases.Service
	rpr model.ProcessRuleRepository
	rq  *orderqueueusecases.Service
}

func NewService(ro model.OrderRepository) *Service {
	return &Service{ro: ro}
}

func (s *Service) AddDependencies(rs model.ShiftRepository, rgi *groupitemusecases.Service, rp *orderprocessusecases.Service, rpr model.ProcessRuleRepository, rq *orderqueueusecases.Service) {
	s.rs = rs
	s.rgi = rgi
	s.rp = rp
	s.rpr = rpr
	s.rq = rq
}
