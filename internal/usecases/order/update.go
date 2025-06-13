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
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

func (s *OrderService) PendingOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()

	processRules, err := s.rpr.GetMapProcessRulesByFirstOrder(ctx)
	if err != nil {
		return err
	}

	groupItemIDs := []uuid.UUID{}
	for i, groupItem := range order.GroupItems {
		if groupItem.Status != orderentity.StatusGroupStaging {
			continue
		}

		if !groupItem.UseProcessRule {
			order.GroupItems[i].PendingGroupItem()
			order.GroupItems[i].StartGroupItem()
			order.GroupItems[i].ReadyGroupItem()
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
			OrderNumber:   order.OrderNumber,
			OrderType:     orderprocessentity.GetTypeOrderProcessFromOrder(order.OrderType),
			GroupItemID:   groupItem.ID,
			ProcessRuleID: processRuleID,
		}

		// Create process for each group item
		if _, err := s.sop.CreateProcess(ctx, createProcessInput); err != nil {
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

		if _, err := s.sq.StartQueue(ctx, startQueueInput); err != nil {
			return err
		}
	}

	orderModel.FromDomain(order)
	if err := s.ro.PendingOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) ReadyOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	if err = order.ReadyOrder(); err != nil {
		return err
	}

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	if order.Delivery != nil {
		dtoDelivery := &entitydto.IDRequest{
			ID: order.Delivery.ID,
		}
		if err := s.sd.ReadyOrderDelivery(ctx, dtoDelivery); err != nil {
			return err
		}

	} else if order.Pickup != nil {
		dtoPickup := &entitydto.IDRequest{
			ID: order.Pickup.ID,
		}
		if err := s.sp.ReadyOrder(ctx, dtoPickup); err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderService) FinishOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	if err = order.FinishOrder(); err != nil {
		return err
	}

	// Update order to new shift
	currentShift, _ := s.rs.GetCurrentShift(ctx)
	if currentShift == nil {
		return fmt.Errorf("must open a new shift")
	}

	order.ShiftID = currentShift.ID

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, dtoOrderID *entitydto.IDRequest) (err error) {
	orderModel, err := s.ro.GetOrderById(ctx, dtoOrderID.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	if err = order.CancelOrder(); err != nil {
		return err
	}

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	if order.Delivery != nil {
		deliveryDtoID := entitydto.NewIdRequest(order.Delivery.ID)
		if err := s.sd.CancelOrderDelivery(ctx, deliveryDtoID); err != nil {
			return err
		}
	} else if order.Pickup != nil {
		pickupDtoID := entitydto.NewIdRequest(order.Pickup.ID)
		if err := s.sp.CancelOrderPickup(ctx, pickupDtoID); err != nil {
			return err
		}
	} else if order.Table != nil {
		tableDtoID := entitydto.NewIdRequest(order.Table.ID)
		if err := s.st.CancelOrderTable(ctx, tableDtoID); err != nil {
			return err
		}
	}

	reason := "order canceled"

	for _, groupItem := range order.GroupItems {
		dtoGroupItemID := entitydto.NewIdRequest(groupItem.ID)
		if err = s.sgi.CancelGroupItem(ctx, dtoGroupItemID); err != nil {
			return err
		}

		processes, err := s.sop.GetProcessesByGroupItemID(ctx, dtoGroupItemID)
		if err != nil {
			return err
		}

		if len(processes) == 0 {
			continue
		}

		for _, process := range processes {
			dtoProcessID := entitydto.NewIdRequest(process.ID)
			orderProcessCancelDTO := &orderprocessdto.OrderProcessCancelDTO{Reason: &reason}
			if err = s.sop.CancelProcess(ctx, dtoProcessID, orderProcessCancelDTO); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *OrderService) ArchiveOrder(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	if err = order.ArchiveOrder(); err != nil {
		return err
	}

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) UnarchiveOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	if err = order.UnarchiveOrder(); err != nil {
		return err
	}

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) AddPayment(ctx context.Context, dto *entitydto.IDRequest, dtoPayment *orderdto.OrderPaymentCreateDTO) error {
	orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	if err = order.ValidatePayments(); err != nil {
		return err
	}

	paymentOrder, err := dtoPayment.ToDomain(order)
	if err != nil {
		return err
	}

	order.AddPayment(paymentOrder)

	order.CalculateTotalPrice()

	paymentOrderModel := &model.PaymentOrder{}
	paymentOrderModel.FromDomain(paymentOrder)
	if err := s.ro.AddPaymentOrder(ctx, paymentOrderModel); err != nil {
		return err
	}

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) UpdateOrderObservation(ctx context.Context, dtoId *entitydto.IDRequest, dto *orderdto.OrderUpdateObservationDTO) error {
	orderModel, err := s.ro.GetOrderById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	dto.UpdateDomain(order)

	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) UpdateOrderTotal(ctx context.Context, id string) error {
	orderModel, err := s.ro.GetOrderById(ctx, id)
	if err != nil {
		return err
	}

	order := orderModel.ToDomain()

	order.CalculateTotalPrice()

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}
