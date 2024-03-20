package itemusecases

import (
	"context"

	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) StartItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = item.StartItem(); err != nil {
		return err
	}

	for i := range item.AdditionalItems {
		if err = item.AdditionalItems[i].StartItem(); err != nil {
			return err
		}
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), false)

	if err != nil {
		return err
	}

	if groupItem.Status == groupitementity.StatusGroupPending {
		groupItem.StartGroupItem()
	}

	if err = s.ri.UpdateItem(ctx, item); err != nil {
		return err
	}

	for i := range item.AdditionalItems {
		if err = s.ri.UpdateItem(ctx, &item.AdditionalItems[i]); err != nil {
			return err
		}
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) ReadyItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = item.ReadyItem(); err != nil {
		return err
	}

	for i := range item.AdditionalItems {
		if err = item.AdditionalItems[i].ReadyItem(); err != nil {
			return err
		}
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), false)

	if err != nil {
		return err
	}

	if err = s.ri.UpdateItem(ctx, item); err != nil {
		return err
	}

	for i := range item.AdditionalItems {
		if err = s.ri.UpdateItem(ctx, &item.AdditionalItems[i]); err != nil {
			return err
		}
	}

	isAllItemsReady := true
	for i := range groupItem.Items {
		if groupItem.Items[i].Status != itementity.StatusItemReady {
			isAllItemsReady = false
			break
		}
	}

	if !isAllItemsReady {
		return nil
	}

	if err = groupItem.ReadyGroupItem(); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) CancelItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	item.CancelItem()

	if err = s.ri.UpdateItem(ctx, item); err != nil {
		return err
	}

	for i := range item.AdditionalItems {
		item.AdditionalItems[i].CancelItem()
		if err = s.ri.UpdateItem(ctx, &item.AdditionalItems[i]); err != nil {
			return err
		}
	}

	return nil
}
