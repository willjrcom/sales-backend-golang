package ordertableusecases

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

type Service struct {
	rto orderentity.OrderTableRepository
	rt  tableentity.TableRepository
	os  *orderusecases.Service
}

func NewService(rto orderentity.OrderTableRepository, rt tableentity.TableRepository, os *orderusecases.Service) *Service {
	return &Service{rto: rto, rt: rt, os: os}
}
