package orderusecases

import (
	"context"

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

func (s *GroupItemService) CancelGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	groupItem.CancelGroupItem()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItem(ctx, groupItemModel); err != nil {
		return err
	}

	if err := s.UpdateGroupItemTotal(ctx, groupItemModel.ID.String()); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	reason := "group item canceled"
	processes, err := s.sop.GetProcessesByGroupItemID(ctx, dto)
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		return nil
	}

	for _, process := range processes {
		dtoProcessID := entitydto.NewIdRequest(process.ID)
		orderProcessCancelDTO := &orderprocessdto.OrderProcessCancelDTO{Reason: &reason}
		if err = s.sop.CancelProcess(ctx, dtoProcessID, orderProcessCancelDTO); err != nil {
			return err
		}
	}

	return nil
}
