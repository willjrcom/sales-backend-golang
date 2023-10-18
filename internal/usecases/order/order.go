package orderusecases

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Service struct {
	ro orderentity.OrderRepository
}

func NewService(ro orderentity.OrderRepository) *Service {
	return &Service{ro: ro}
}
