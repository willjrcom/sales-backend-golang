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

	for i := range groupItem.Items {
		if err = groupItem.Items[i].StartItem(); err != nil {
			return err
		}

		for j := range groupItem.Items[i].AdditionalItems {
			if err = groupItem.Items[i].AdditionalItems[j].StartItem(); err != nil {
				return err
			}
		}
	}

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	for i := range groupItem.Items {
		if err = s.ri.UpdateItem(ctx, &groupItem.Items[i]); err != nil {
			return err
		}

		for j := range groupItem.Items[i].AdditionalItems {
			if err = s.ri.UpdateItem(ctx, &groupItem.Items[i].AdditionalItems[j]); err != nil {
				return err
			}
		}
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

	for i := range groupItem.Items {
		if err = groupItem.Items[i].ReadyItem(); err != nil {
			return err
		}

		for j := range groupItem.Items[i].AdditionalItems {
			if err = groupItem.Items[i].AdditionalItems[j].ReadyItem(); err != nil {
				return err
			}
		}
	}

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	for i := range groupItem.Items {
		if err = s.ri.UpdateItem(ctx, &groupItem.Items[i]); err != nil {
			return err
		}

		for j := range groupItem.Items[i].AdditionalItems {
			if err = s.ri.UpdateItem(ctx, &groupItem.Items[i].AdditionalItems[j]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) CancelGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	groupItem.CancelGroupItem()

	for i := range groupItem.Items {
		groupItem.Items[i].CancelItem()

		for j := range groupItem.Items[i].AdditionalItems {
			groupItem.Items[i].AdditionalItems[j].CancelItem()
		}
	}

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	for i := range groupItem.Items {
		if err = s.ri.UpdateItem(ctx, &groupItem.Items[i]); err != nil {
			return err
		}

		for j := range groupItem.Items[i].AdditionalItems {
			if err = s.ri.UpdateItem(ctx, &groupItem.Items[i].AdditionalItems[j]); err != nil {
				return err
			}
		}
	}

	return nil
}
