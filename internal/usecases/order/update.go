package orderusecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	orderprocessdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_process"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
)

func (s *OrderService) PendingOrder(ctx context.Context, dto *entitydto.IDRequest) error {
	orderModel, err := s.ro.GetOrderById(ctx, dto.ID.String())
	if err != nil {
		return err
	}

	order := orderModel.ToDomain()

	for _, groupItem := range order.GroupItems {
		if math.Mod(groupItem.Quantity, 1) != 0 {
			return orderentity.ErrQuantityNotInteger
		}
	}

	processRules, err := s.rpr.GetMapProcessRulesByFirstOrder(ctx)
	if err != nil {
		return err
	}

	company, err := s.sc.GetCompany(ctx)
	if err != nil {
		return err
	}

	for i, groupItem := range order.GroupItems {
		if groupItem.Status != orderentity.StatusGroupStaging {
			continue
		}

		if !groupItem.UseProcessRule {
			order.GroupItems[i].PendingGroupItem()
			order.GroupItems[i].StartGroupItem()
			order.GroupItems[i].ReadyGroupItem()
		}

		// Generate snapshot for all group items (including those that move directly to Started/Ready)
		groupItemDTO := &groupitemdto.GroupItemDTO{}
		groupItemDTO.FromDomain(&order.GroupItems[i])

		if data, err := json.Marshal(groupItemDTO); err == nil {
			snapshot := &model.OrderGroupItemSnapshot{
				Entity:      entitymodel.FromDomain(entity.NewEntity()),
				GroupItemID: groupItem.ID,
				Data:        data,
			}

			if err := s.sgi.r.UpsertGroupItemSnapshot(ctx, snapshot); err != nil {
				fmt.Printf("error upserting group item snapshot: %v\n", err)
			}
		}

		if !groupItem.UseProcessRule {
			continue
		}

		if groupItem.Category.NeedPrint && s.rabbitmq != nil {
			path := rabbitmq.GROUP_ITEM_PATH + groupItem.ID.String()
			if err := s.rabbitmq.SendPrintMessage(rabbitmq.GROUP_ITEM_EX, company.SchemaName, path, groupItem.Category.PrinterName); err != nil {
				fmt.Println("error sending message to rabbitmq: " + err.Error())
			}
		}

		processRuleID, ok := processRules[groupItem.CategoryID]
		if !ok {
			fmt.Println("process rule not found for category ID: " + groupItem.CategoryID.String())
			continue
		}

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
	for _, groupItem := range order.GroupItems {
		if groupItem.Status != orderentity.StatusGroupStaging && !groupItem.UseProcessRule {
			continue
		}

		startQueueInput := &orderqueuedto.QueueCreateDTO{
			GroupItemID: groupItem.ID,
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

	if EnablePrintOrderOnShipOrder, _ := company.Preferences.GetBool(companyentity.EnablePrintOrderOnShipOrder); EnablePrintOrderOnShipOrder {
		printerName, _ := company.Preferences.GetString(companyentity.PrinterOrder)
		if s.rabbitmq != nil {
			path := rabbitmq.ORDER_PATH + order.ID.String()
			if err := s.rabbitmq.SendPrintMessage(rabbitmq.ORDER_EX, company.SchemaName, path, printerName); err != nil {
				fmt.Println("error sending message to rabbitmq: " + err.Error())
			}
		}
	}

	return nil
}

// restoreStockFromOrder restaura estoque dos produtos do pedido cancelado
func (s *OrderService) restoreStockFromOrder(ctx context.Context, order *orderentity.Order) error {
	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	employee, err := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
	if err != nil {
		return err
	}

	for _, groupItem := range order.GroupItems {
		s.sgi.restoreStockFromGroupItem(ctx, &groupItem, employee.ID)
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

func (s *OrderService) CancelOrder(ctx context.Context, dtoOrderID *entitydto.IDRequest, isSystem bool) (err error) {
	orderModel, err := s.ro.GetOrderById(ctx, dtoOrderID.ID.String())

	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	if err = order.CancelOrder(); err != nil {
		return err
	}

	// Update order to new shift
	if !isSystem {
		currentShift, _ := s.rs.GetCurrentShift(ctx)
		if currentShift == nil {
			return fmt.Errorf("must open a new shift")
		}

		order.ShiftID = currentShift.ID
	}

	// Restaurar estoque dos produtos do pedido cancelado
	if err := s.restoreStockFromOrder(ctx, order); err != nil {
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

	reason := "order cancelled"
	cancelDTO := &groupitemdto.OrderGroupItemCancelDTO{Reason: &reason}

	for _, groupItem := range order.GroupItems {
		if err = s.sgi.CancelGroupItem(ctx, groupItem.ID.String(), cancelDTO); err != nil {
			return err
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

	order.CalculateTotalPaid()

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

func (s *OrderService) UpdateTotalFees(ctx context.Context, id string) error {
	orderModel, err := s.ro.GetOnlyOrderById(ctx, id)
	if err != nil {
		return err
	}

	order := orderModel.ToDomain()
	order.CalculateFees()
	order.CalculateTotal()

	orderModel.FromDomain(order)
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

	company, err := s.sc.GetCompany(ctx)
	if err != nil {
		return err
	}

	if order.Delivery != nil {
		minOrderForFree, _ := company.Preferences.GetDecimal(companyentity.MinOrderValueForFreeDelivery)

		order.Delivery.IsDeliveryFree = false
		if !minOrderForFree.IsZero() && order.SubTotal.GreaterThan(minOrderForFree) {
			order.Delivery.IsDeliveryFree = true
		}

		deliveryOrder := &model.OrderDelivery{}
		deliveryOrder.FromDomain(order.Delivery)
		if err := s.rdo.UpdateOrderDelivery(ctx, deliveryOrder); err != nil {
			return err
		}
	}

	order.CalculateTotalOrder()

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}
