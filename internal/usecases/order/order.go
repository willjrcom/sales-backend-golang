package orderusecases

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
	orderprocessusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_process"
	queueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/queue"
)

type Service struct {
	ro  orderentity.OrderRepository
	rs  shiftentity.ShiftRepository
	rgi *groupitemusecases.Service
	rp  *orderprocessusecases.Service
	rpr productentity.ProcessRuleRepository
	rq  *queueusecases.Service
}

func NewService(ro orderentity.OrderRepository, rs shiftentity.ShiftRepository, rgi *groupitemusecases.Service, rp *orderprocessusecases.Service, rpr productentity.ProcessRuleRepository, rq *queueusecases.Service) *Service {
	return &Service{ro: ro, rs: rs, rgi: rgi, rp: rp, rpr: rpr, rq: rq}
}
