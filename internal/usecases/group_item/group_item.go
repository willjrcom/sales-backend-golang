package groupitemusecases

import (
	"context"
	"errors"

	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
)

var (
	ErrItemsFinished = errors.New("items already finished")
)

type Service struct {
	ri  itementity.ItemRepository
	rgi groupitementity.GroupItemRepository
}

func NewService(ri itementity.ItemRepository, rgi groupitementity.GroupItemRepository) *Service {
	return &Service{ri: ri, rgi: rgi}
}

func (s *Service) GetGroupByID(ctx context.Context, dto *entitydto.IdRequest) (groupItem *groupitementity.GroupItem, err error) {
	return s.rgi.GetGroupByID(ctx, dto.ID.String(), true)
}

func (s *Service) GetAllGroupsByStatus(ctx context.Context, dto *groupitemdto.GroupStatusInput) (groups []groupitementity.GroupItem, err error) {
	return s.rgi.GetAllGroupsByStatus(ctx, dto.Status)
}

func (s *Service) StartGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if err = groupItem.StartGroupItem(); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) ReadyGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if err = groupItem.ReadyGroupItem(); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) CancelGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	groupItem.CancelGroupItem()

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) FinishGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if err = groupItem.FinishGroupItem(); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) DeleteGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if err = groupItem.DeleteGroupItem(); err != nil {
		return err
	}

	return s.rgi.DeleteGroupItem(ctx, groupItem.ID.String())
}
