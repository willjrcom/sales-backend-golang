package groupitemusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) StartGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if err = groupItem.StartGroupItem(); err != nil {
		return err
	}

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	return nil
}

func (s *Service) ReadyGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if err = groupItem.ReadyGroupItem(); err != nil {
		return err
	}

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	groupItem.CancelGroupItem()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	return nil
}
