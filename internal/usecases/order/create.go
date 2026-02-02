package orderusecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrSubscriptionExpired         = errors.New("assinatura expirada: realize o pagamento da mensalidade para continuar")
	ErrCompanyBlocked              = errors.New("conta bloqueada: entre em contato com o suporte")
	ErrCompanySubscriptionNotFound = errors.New("assinatura não encontrada")
)

func (s *OrderService) ValidateSubscription(ctx context.Context) error {
	companyModel, err := s.sc.GetCompany(ctx)
	if err != nil {
		return err
	}

	if companyModel.IsBlocked {
		return ErrCompanyBlocked
	}

	activeSub, _ := s.companySubscriptionRepo.GetActiveSubscription(ctx, companyModel.ID)
	if activeSub == nil {
		return ErrCompanySubscriptionNotFound
	}

	if activeSub.EndDate.Before(time.Now().UTC()) {
		return ErrSubscriptionExpired
	}

	return nil
}

func (s *OrderService) CreateDefaultOrder(ctx context.Context) (uuid.UUID, error) {
	if s.sc != nil {
		if err := s.ValidateSubscription(ctx); err != nil {
			return uuid.Nil, err
		}
	}

	shiftModel, err := s.rs.GetCurrentShift(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("must open a new shift")
	}

	currentOrderNumber, err := s.rs.IncrementCurrentOrder(ctx, shiftModel.ID.String())
	if err != nil {
		return uuid.Nil, err
	}

	order := orderentity.NewDefaultOrder(shiftModel.ID, currentOrderNumber, shiftModel.AttendantID)

	orderModel := &model.Order{}
	orderModel.FromDomain(order)

	if err := s.ro.CreateOrder(ctx, orderModel); err != nil {
		return uuid.Nil, err
	}

	return order.ID, nil
}
