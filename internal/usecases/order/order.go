package orderusecases

import (
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrOrderMustBePending = errors.New("order must be pending")
)

type Service struct {
	ro orderentity.OrderRepository
}

func NewService(ro orderentity.OrderRepository) *Service {
	return &Service{ro: ro}
}
