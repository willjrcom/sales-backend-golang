package groupitemusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
)

func (s *Service) UpdateScheduleGroupItem(ctx context.Context, dtoId *entitydto.IdRequest, dto *groupitemdto.UpdateScheduleGroupItem) (err error) {
	startAt, err := dto.ToModel()

	if err != nil {
		return err
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, dtoId.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem.Schedule(startAt)

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) StartGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	if err = groupItem.StartGroupItem(); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) ReadyGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	if err = groupItem.ReadyGroupItem(); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) CancelGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	groupItem.CancelGroupItem()

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}
