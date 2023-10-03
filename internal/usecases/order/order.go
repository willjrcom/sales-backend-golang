package orderusecases

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type Service struct {
	ro orderentity.Repository
	ra addressentity.Repository
}

func NewService(ro orderentity.Repository, ra addressentity.Repository) *Service {
	return &Service{ro: ro, ra: ra}
}
