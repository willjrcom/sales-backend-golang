package ordertableusecases

import (
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
)

type Service struct {
	rto model.OrderTableRepository
	rt  model.TableRepository
	os  *orderusecases.Service
}

func NewService(rto model.OrderTableRepository) *Service {
	return &Service{rto: rto}
}

func (s *Service) AddDependencies(rt model.TableRepository, os *orderusecases.Service) {
	s.rt = rt
	s.os = os
}
