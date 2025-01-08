package orderusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrItemsFinished              = errors.New("items already finished")
	ErrSizeMustBeTheSame          = errors.New("size must be the same")
	ErrComplementItemAlreadyAdded = errors.New("complement item already added")
	ErrComplementItemNotFound     = errors.New("complement item not found")
)

func (s *GroupItemService) GetGroupByID(ctx context.Context, dto *entitydto.IDRequest) (groupItemDTO *groupitemdto.GroupItemDTO, err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return nil, err
	}

	groupItem := groupItemModel.ToDomain()
	groupItemDTO = &groupitemdto.GroupItemDTO{}
	groupItemDTO.FromDomain(groupItem)
	return
}

func (s *GroupItemService) DeleteGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItem, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	var complementItemID *string
	if groupItem.ComplementItemID != nil {
		complementItemID = new(string)
		*complementItemID = groupItem.ComplementItemID.String()
	}

	if err := s.r.DeleteGroupItem(ctx, groupItem.ID.String(), complementItemID); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return err
	}

	return nil
}

func (s *GroupItemService) AddComplementItem(ctx context.Context, dto *entitydto.IDRequest, dtoComplement *entitydto.IDRequest) (err error) {
	groupItem, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if groupItem.ComplementItemID != nil {
		return ErrComplementItemAlreadyAdded
	}

	productComplement, err := s.rp.GetProductById(ctx, dtoComplement.ID.String())

	if err != nil {
		return err
	}

	if groupItem.Size != productComplement.Size.Name {
		return ErrSizeMustBeTheSame
	}

	// Is valid complement to this category
	found := false
	for _, complementCategory := range groupItem.Category.ComplementCategories {
		if productComplement.Category.ID == complementCategory.ID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("complement category does not belong to this category")
	}

	itemComplement := orderentity.NewItem(productComplement.Name, productComplement.Price, groupItem.Quantity, groupItem.Size, productComplement.ID, productComplement.CategoryID)

	itemComplementModel := &model.Item{}
	itemComplementModel.FromDomain(itemComplement)
	if err = s.ri.AddItem(ctx, itemComplementModel); err != nil {
		return err
	}

	groupItem.ComplementItemID = &itemComplement.ID

	if err := s.r.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return err
	}

	return nil
}

func (s *GroupItemService) DeleteComplementItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItem, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if groupItem.ComplementItemID == nil {
		return ErrComplementItemNotFound
	}

	if err := s.ri.DeleteItem(ctx, groupItem.ComplementItemID.String()); err != nil {
		return err
	}

	groupItem.ComplementItemID = nil
	groupItem.ComplementItem = nil

	if err := s.r.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return err
	}

	return nil
}
