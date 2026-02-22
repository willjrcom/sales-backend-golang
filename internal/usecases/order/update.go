package orderusecases

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	orderprocessdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_process"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
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

	// Controle de estoque - debitar estoque dos produtos
	if err := s.debitStockFromOrder(ctx, order); err != nil {
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
			continue
		}

		if groupItem.Category.NeedPrint {
			schemaName, err := database.GetCurrentSchema(ctx)
			if err != nil {
				fmt.Println("error getting schema name: " + err.Error())
			} else if s.rabbitmq != nil {
				if err := s.rabbitmq.SendMessage(schemaName, rabbitmq.GROUP_ITEM_RK, groupItem.ID.String()); err != nil {
					fmt.Println("error sending message to rabbitmq: " + err.Error())
				}
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

	return nil
}

// debitStockFromOrder debita estoque dos produtos do pedido
func (s *OrderService) debitStockFromOrder(ctx context.Context, order *orderentity.Order) error {
	for _, groupItem := range order.GroupItems {
		if groupItem.Status != orderentity.StatusGroupStaging {
			continue
		}

		for _, item := range groupItem.Items {
			if item.ProductID != uuid.Nil {
				fmt.Printf("DEBUG: Produto %s - Quantidade: %f\n", item.Name, item.Quantity)

				// Buscar estoque do produto/variação
				stockModel, err := s.stockRepo.GetStockByVariationID(ctx, item.ProductVariationID.String())
				if err != nil {
					// Fallback para buscar apenas por ProductID se não houver variação específica (ex: adicionais sem tamanho)
					stocks, err := s.stockRepo.GetStockByProductID(ctx, item.ProductID.String())
					if err != nil || len(stocks) == 0 {
						// Se não há controle de estoque para o produto, continuar
						fmt.Printf("Produto/Variação %s não tem controle de estoque configurado\n", item.Name)
						continue
					}
					stockModel = &stocks[0]
				}

				stock := stockModel.ToDomain()

				// Reservar estoque (permite estoque negativo)
				movement, err := stock.ReserveStock(
					decimal.NewFromFloat(item.Quantity),
					order.ID,
					*order.AttendantID,
					item.Price,
					item.TotalPrice,
				)
				if err != nil {
					fmt.Printf("Erro ao reservar estoque para produto %s: %v\n", item.Name, err)
					continue
				}

				// Salvar movimento
				movementModel := &model.StockMovement{}
				movementModel.FromDomain(movement)
				if err := s.stockMovementRepo.CreateMovement(ctx, movementModel); err != nil {
					fmt.Printf("Erro ao salvar movimento de estoque: %v\n", err)
					continue
				}

				// Atualizar estoque
				stockModel.FromDomain(stock)
				if err := s.stockRepo.UpdateStock(ctx, stockModel); err != nil {
					fmt.Printf("Erro ao atualizar estoque: %v\n", err)
					continue
				}

				fmt.Printf("Estoque debitado para produto %s: %f\n", item.Name, item.Quantity)
			}
		}
	}

	return nil
}

// restoreStockFromOrder restaura estoque dos produtos do pedido cancelado
func (s *OrderService) restoreStockFromOrder(ctx context.Context, order *orderentity.Order) error {
	for _, groupItem := range order.GroupItems {
		for _, item := range groupItem.Items {
			if item.ProductID != uuid.Nil {
				fmt.Printf("DEBUG: Restaurando estoque para produto %s - Quantidade: %f\n", item.Name, item.Quantity)

				// Buscar estoque do produto/variação
				stockModel, err := s.stockRepo.GetStockByVariationID(ctx, item.ProductVariationID.String())
				if err != nil {
					// Fallback para buscar apenas por ProductID
					stocks, err := s.stockRepo.GetStockByProductID(ctx, item.ProductID.String())
					if err != nil || len(stocks) == 0 {
						// Se não há controle de estoque para o produto, continuar
						fmt.Printf("Produto/Variação %s não tem controle de estoque configurado\n", item.Name)
						continue
					}
					stockModel = &stocks[0]
				}

				stock := stockModel.ToDomain()

				// Restaurar estoque
				movement, err := stock.RestoreStock(
					decimal.NewFromFloat(item.Quantity),
					order.ID,
					*order.AttendantID,
					item.Price,
					item.TotalPrice,
				)
				if err != nil {
					fmt.Printf("Erro ao restaurar estoque para produto %s: %v\n", item.Name, err)
					continue
				}

				// Salvar movimento
				movementModel := &model.StockMovement{}
				movementModel.FromDomain(movement)
				if err := s.stockMovementRepo.CreateMovement(ctx, movementModel); err != nil {
					fmt.Printf("Erro ao salvar movimento de estoque: %v\n", err)
					continue
				}

				// Atualizar estoque
				stockModel.FromDomain(stock)
				if err := s.stockRepo.UpdateStock(ctx, stockModel); err != nil {
					fmt.Printf("Erro ao atualizar estoque: %v\n", err)
					continue
				}

				fmt.Printf("Estoque restaurado para produto %s: %f\n", item.Name, item.Quantity)
			}
		}
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

	for _, groupItem := range order.GroupItems {
		dtoGroupItemID := entitydto.NewIdRequest(groupItem.ID)
		if err = s.sgi.CancelGroupItem(ctx, dtoGroupItemID); err != nil {
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

	if err := s.updateFreeDelivery(ctx, order); err != nil {
		return err
	}

	orderModel.FromDomain(order)
	if err := s.ro.UpdateOrder(ctx, orderModel); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) updateFreeDelivery(ctx context.Context, order *orderentity.Order) error {
	if order.Delivery == nil {
		return nil
	}

	IsDeliveryFreeUpdated := order.Delivery.IsDeliveryFree

	company, err := s.sc.GetCompany(ctx)
	if err != nil {
		return err
	}

	isMinOrderForFreeEnabled, err := company.Preferences.GetBool("enable_min_order_value_for_free_delivery")
	if err != nil {
		return err
	}

	if isMinOrderForFreeEnabled {
		minOrderForFree, err := company.Preferences.GetDecimal("min_order_value_for_free_delivery")
		if err != nil {
			return err
		}

		order.Delivery.IsDeliveryFree = false
		if order.TotalPayable.GreaterThan(minOrderForFree) {
			order.Delivery.IsDeliveryFree = true
		}
	}

	if IsDeliveryFreeUpdated == order.Delivery.IsDeliveryFree {
		return nil
	}

	// Run again with IsDeliveryFree changes
	order.CalculateTotalPrice()

	deliveryModel := &model.OrderDelivery{}

	deliveryModel.FromDomain(order.Delivery)
	if err := s.rdo.UpdateOrderDelivery(ctx, deliveryModel); err != nil {
		return err
	}

	return nil
}
