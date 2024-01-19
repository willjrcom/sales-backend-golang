package orderusecases

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

var (
	ErrOrderMustBePending = errors.New("order must be pending")
)

type Service struct {
	ro orderentity.OrderRepository
	rs shiftentity.ShiftRepository
}

func NewService(ro orderentity.OrderRepository, rs shiftentity.ShiftRepository) *Service {
	return &Service{ro: ro, rs: rs}
}
