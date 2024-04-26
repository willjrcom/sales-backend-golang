package orderusecases

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
)

type Service struct {
	ro  orderentity.OrderRepository
	rs  shiftentity.ShiftRepository
	rgi *groupitemusecases.Service
}

func NewService(ro orderentity.OrderRepository, rs shiftentity.ShiftRepository, rgi *groupitemusecases.Service) *Service {
	return &Service{ro: ro, rs: rs, rgi: rgi}
}
