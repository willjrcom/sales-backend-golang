package tableorderusecases

import (
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type Service struct {
	rto orderentity.TableOrderRepository
	rt  tableentity.TableRepository
}

func NewService(rto orderentity.TableOrderRepository, rt tableentity.TableRepository) *Service {
	return &Service{rto: rto, rt: rt}
}
