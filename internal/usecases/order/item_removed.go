package orderusecases

import (
	"context"
	"errors"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
)

func (s *ItemService) AddRemovedItem(ctx context.Context, dtoID *entitydto.IDRequest, dto *itemdto.RemovedItemDTO) (err error) {
	name, err := dto.ToDomain()
	if err != nil {
		return err
	}

	itemModel, err := s.ri.GetItemById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	item := itemModel.ToDomain()

	category, err := s.rc.GetCategoryById(ctx, item.CategoryID.String())
	if err != nil {
		return err
	}

	found := false
	for _, removedItem := range category.RemovableIngredients {
		if removedItem == *name {
			found = true
		}
	}

	if !found {
		return errors.New("removed item not found on category")
	}

	// Already added
	for _, removedItem := range item.RemovedItems {
		if removedItem == *name {
			return nil
		}
	}

	item.AddRemovedItem(*name)

	itemModel.FromDomain(item)
	return s.ri.UpdateItem(ctx, itemModel)
}

func (s *ItemService) RemoveRemovedItem(ctx context.Context, dtoID *entitydto.IDRequest, dto *itemdto.RemovedItemDTO) (err error) {
	name, err := dto.ToDomain()
	if err != nil {
		return err
	}

	itemModel, err := s.ri.GetItemById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	item := itemModel.ToDomain()
	item.RemoveRemovedItem(*name)

	itemModel.FromDomain(item)
	return s.ri.UpdateItem(ctx, itemModel)
}
