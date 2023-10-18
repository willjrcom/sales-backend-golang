package tableorderusecases

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Service struct {
	rt orderentity.TableRepository
}

func NewService(rt orderentity.TableRepository) *Service {
	return &Service{rt: rt}
}
