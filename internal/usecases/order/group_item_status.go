package orderusecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	orderprocessdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_process"
)

func (s *GroupItemService) UpdateScheduleGroupItem(ctx context.Context, dtoId *entitydto.IDRequest, dto *groupitemdto.OrderGroupItemUpdateScheduleDTO) (err error) {
	startAt, err := dto.ToDomain()

	if err != nil {
		return err
	}

	groupItemModel, err := s.r.GetGroupByID(ctx, dtoId.ID.String(), false)
	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	groupItem.Schedule(startAt)

	groupItemModel.FromDomain(groupItem)
	return s.r.UpdateGroupItem(ctx, groupItemModel)
}

func (s *GroupItemService) UpdateObservationGroupItem(ctx context.Context, dtoId *entitydto.IDRequest, dto *groupitemdto.OrderGroupItemUpdateObservationDTO) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dtoId.ID.String(), false)
	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	groupItem.Observation = dto.Observation

	groupItemModel.FromDomain(groupItem)
	return s.r.UpdateGroupItem(ctx, groupItemModel)
}

func (s *GroupItemService) StartGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	if err = groupItem.StartGroupItem(); err != nil {
		return err
	}

	groupItemModel.FromDomain(groupItem)
	return s.r.UpdateGroupItem(ctx, groupItemModel)
}

func (s *GroupItemService) ReadyGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	if err = groupItem.ReadyGroupItem(); err != nil {
		return err
	}

	groupItemModel.FromDomain(groupItem)
	return s.r.UpdateGroupItem(ctx, groupItemModel)
}

func (s *GroupItemService) CancelGroupItem(ctx context.Context, id string, dto *groupitemdto.OrderGroupItemCancelDTO) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, id, true)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	// Restaurar estoque dos itens do grupo se ainda NÃO estiver cancelado
	if groupItem.Status != orderentity.StatusGroupCancelled {
		userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
		if !ok {
			return errors.New("context user not found")
		}

		userIDUUID := uuid.MustParse(userID)
		employee, err := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
		if err != nil {
			return err
		}

		if err := s.restoreStockFromGroupItemWithTx(ctx, tx, groupItem, employee.ID); err != nil {
			return err
		}
	}

	if dto.Reason == nil || *dto.Reason == "" {
		return fmt.Errorf("reason is required")
	}

	groupItem.CancelledReason = *dto.Reason
	groupItem.CancelGroupItem()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItemWithTx(ctx, tx, groupItemModel); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	dtoId := &entitydto.IDRequest{ID: groupItemModel.ID}
	processes, err := s.sop.GetProcessesByGroupItemID(ctx, dtoId)
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		return nil
	}

	for _, process := range processes {
		dtoProcessID := entitydto.NewIdRequest(process.ID)
		orderProcessCancelDTO := &orderprocessdto.OrderProcessCancelDTO{Reason: dto.Reason}
		if err = s.sop.CancelProcess(ctx, dtoProcessID, orderProcessCancelDTO); err != nil {
			return err
		}
	}

	return nil
}

func (s *GroupItemService) restoreStockFromGroupItem(ctx context.Context, groupItem *orderentity.GroupItem, attendantID uuid.UUID) error {
	for _, item := range groupItem.Items {
		if item.ProductID != uuid.Nil {
			if err := s.si.RestoreStockFromItem(ctx, &item, groupItem, attendantID); err != nil {
				fmt.Printf("Aviso: erro ao restaurar estoque do item %s: %v\n", item.Name, err)
			}
		}

		for _, addItem := range item.AdditionalItems {
			if addItem.ProductID != uuid.Nil {
				if err := s.si.RestoreStockFromItem(ctx, &addItem, groupItem, attendantID); err != nil {
					fmt.Printf("Aviso: erro ao restaurar estoque do adicional %s: %v\n", addItem.Name, err)
				}
			}
		}
	}

	if groupItem.ComplementItem != nil && groupItem.ComplementItem.ProductID != uuid.Nil {
		if err := s.si.RestoreStockFromItem(ctx, groupItem.ComplementItem, groupItem, attendantID); err != nil {
			fmt.Printf("Aviso: erro ao restaurar estoque do complemento: %v\n", err)
		}
	}

	return nil
}

func (s *GroupItemService) restoreStockFromGroupItemWithTx(ctx context.Context, tx *bun.Tx, groupItem *orderentity.GroupItem, attendantID uuid.UUID) error {
	for _, item := range groupItem.Items {
		if item.ProductID != uuid.Nil {
			if err := s.si.restoreStockFromItemWithTx(ctx, tx, &item, groupItem, attendantID); err != nil {
				return err
			}
		}

		for _, addItem := range item.AdditionalItems {
			if addItem.ProductID != uuid.Nil {
				if err := s.si.restoreStockFromItemWithTx(ctx, tx, &addItem, groupItem, attendantID); err != nil {
					return err
				}
			}
		}
	}

	if groupItem.ComplementItem != nil && groupItem.ComplementItem.ProductID != uuid.Nil {
		if err := s.si.restoreStockFromItemWithTx(ctx, tx, groupItem.ComplementItem, groupItem, attendantID); err != nil {
			return err
		}
	}

	return nil
}

// CancelGroupItemSkipStockRestore cancela o grupo sem restaurar estoque.
// Usado quando o estoque já foi restaurado pela chamada anterior (ex: restoreStockFromOrder em CancelOrder).
func (s *GroupItemService) CancelGroupItemSkipStockRestore(ctx context.Context, id string, dto *groupitemdto.OrderGroupItemCancelDTO) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, id, true)
	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if dto.Reason == nil || *dto.Reason == "" {
		return fmt.Errorf("reason is required")
	}

	groupItem.CancelledReason = *dto.Reason
	groupItem.CancelGroupItem()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItemWithTx(ctx, tx, groupItemModel); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	dtoId := &entitydto.IDRequest{ID: groupItemModel.ID}
	processes, err := s.sop.GetProcessesByGroupItemID(ctx, dtoId)
	if err != nil {
		return err
	}

	for _, process := range processes {
		dtoProcessID := entitydto.NewIdRequest(process.ID)
		orderProcessCancelDTO := &orderprocessdto.OrderProcessCancelDTO{Reason: dto.Reason}
		if err = s.sop.CancelProcess(ctx, dtoProcessID, orderProcessCancelDTO); err != nil {
			return err
		}
	}

	return nil
}
