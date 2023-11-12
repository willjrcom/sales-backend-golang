package groupitemusecases

import (
	"context"
	"errors"

	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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
	return s.rgi.GetGroupByID(ctx, dto.ID.String())
}

func (s *Service) DeleteGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String())

	if groupItem.LaunchedAt != nil {
		return ErrItemsFinished
	}

	if len(groupItem.Items) != 0 {
		return
	}

	return s.rgi.DeleteGroupItem(ctx, dto.ID.String())
}

func (s *Service) GetAllPendingGroups(ctx context.Context) (groups []groupitementity.GroupItem, err error) {
	return s.rgi.GetAllPendingGroups(ctx)
}
