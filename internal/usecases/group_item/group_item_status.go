package groupitemusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
)

func (s *Service) UpdateScheduleGroupItem(ctx context.Context, dtoId *entitydto.IDRequest, dto *groupitemdto.OrderGroupItemUpdateScheduleDTO) (err error) {
	startAt, err := dto.ToDomain()

	if err != nil {
		return err
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, dtoId.ID.String(), false)
	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	groupItem.Schedule(startAt)

	groupItemModel.FromDomain(groupItem)
	return s.rgi.UpdateGroupItem(ctx, groupItemModel)
}

func (s *Service) StartGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	if err = groupItem.StartGroupItem(); err != nil {
		return err
	}

	groupItemModel.FromDomain(groupItem)
	return s.rgi.UpdateGroupItem(ctx, groupItemModel)
}

func (s *Service) ReadyGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	if err = groupItem.ReadyGroupItem(); err != nil {
		return err
	}

	groupItemModel.FromDomain(groupItem)
	return s.rgi.UpdateGroupItem(ctx, groupItemModel)
}

func (s *Service) CancelGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()
	groupItem.CancelGroupItem()

	return s.rgi.UpdateGroupItem(ctx, groupItemModel)
}
