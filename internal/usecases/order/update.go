package orderusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	orderprocessdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_process"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
)

func (s *Service) PendingOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	processRules, err := s.rpr.GetMapProcessRulesByFirstOrder(ctx)
	if err != nil {
		return err
	}

	groupItemIDs := []uuid.UUID{}
	for i, groupItem := range order.Groups {
		if groupItem.Status != orderentity.StatusGroupStaging {
			continue
		}

		if !groupItem.UseProcessRule {
			order.Groups[i].PendingGroupItem()
			order.Groups[i].StartGroupItem()
			order.Groups[i].ReadyGroupItem()
			continue
		}

		processRuleID, ok := processRules[groupItem.CategoryID]
		if !ok {
			fmt.Println("process rule not found for category ID: " + groupItem.CategoryID.String())
			continue
		}

		// Append only Staging group items
		groupItemIDs = append(groupItemIDs, groupItem.ID)

		createProcessInput := &orderprocessdto.OrderProcessCreateDTO{
			OrderProcessCommonAttributes: orderprocessentity.OrderProcessCommonAttributes{
				GroupItemID:   groupItem.ID,
				ProcessRuleID: processRuleID,
			},
		}

		// Create process for each group item
		if _, err := s.rp.CreateProcess(ctx, createProcessInput); err != nil {
			return err
		}
	}

	if err = order.PendingOrder(); err != nil {
		return err
	}

	// Create queue for each group item
	for _, groupItemID := range groupItemIDs {
		startQueueInput := &orderqueuedto.QueueCreateDTO{
			GroupItemID: groupItemID,
			JoinedAt:    *order.PendingAt,
		}

		if _, err := s.rq.StartQueue(ctx, startQueueInput); err != nil {
			return err
		}
	}

	if err := s.ro.PendingOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) ReadyOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.ReadyOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) FinishOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.FinishOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelOrder(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.CancelOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	for _, groupItem := range order.Groups {
		dtoID := entitydto.NewIdRequest(groupItem.ID)
		if err = s.rgi.CancelGroupItem(ctx, dtoID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ArchiveOrder(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.ArchiveOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UnarchiveOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.UnarchiveOrder(); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddPayment(ctx context.Context, dto *entitydto.IDRequest, dtoPayment *orderdto.OrderPaymentCreateDTO) error {
	order, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = order.ValidatePayments(); err != nil {
		return err
	}

	paymentOrder, err := dtoPayment.ToDomain(order)
	if err != nil {
		return err
	}

	order.AddPayment(paymentOrder)

	order.CalculateTotalPrice()
	if err := s.ro.AddPaymentOrder(ctx, paymentOrder); err != nil {
		return err
	}

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateOrderObservation(ctx context.Context, dtoId *entitydto.IDRequest, dto *orderdto.OrderUpdateObservationDTO) error {
	order, err := s.ro.GetOrderById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	dto.UpdateDomain(order)

	if err := s.ro.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}
