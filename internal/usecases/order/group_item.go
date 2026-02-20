package orderusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
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

type GroupItemService struct {
	r   model.GroupItemRepository
	ri  model.ItemRepository
	rp  model.ProductRepository
	sop *OrderProcessService
	so  *OrderService
}

func NewGroupItemService(rgi model.GroupItemRepository) *GroupItemService {
	return &GroupItemService{r: rgi}
}

func (s *GroupItemService) AddDependencies(ri model.ItemRepository, rp model.ProductRepository, so *OrderService, sop *OrderProcessService) {
	s.ri = ri
	s.rp = rp
	s.so = so
	s.sop = sop
}

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
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()

	if groupItem.ComplementItemID != nil {
		return ErrComplementItemAlreadyAdded
	}

	productComplementModel, err := s.rp.GetProductById(ctx, dtoComplement.ID.String())

	if err != nil {
		return err
	}

	productComplement := productComplementModel.ToDomain()

	var variation *productentity.ProductVariation
	for _, v := range productComplement.Variations {
		if v.Size != nil && v.Size.Name == groupItem.Size {
			variation = &v
			break
		}
	}

	if variation == nil {
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

	itemComplement := orderentity.NewItem(productComplement.Name, variation.Price, groupItem.Quantity, groupItem.Size, productComplement.ID, productComplement.CategoryID, nil)

	itemComplementModel := &model.Item{}
	itemComplementModel.FromDomain(itemComplement)
	if err = s.ri.AddItem(ctx, itemComplementModel); err != nil {
		return err
	}

	groupItem.ComplementItemID = &itemComplement.ID
	groupItem.ComplementItem = itemComplement

	groupItem.CalculateTotalPrice()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItem(ctx, groupItemModel); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return err
	}

	return nil
}

func (s *GroupItemService) DeleteComplementItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if groupItemModel.ComplementItemID == nil {
		return ErrComplementItemNotFound
	}

	if err := s.ri.DeleteItem(ctx, groupItemModel.ComplementItemID.String()); err != nil {
		return err
	}

	groupItemModel.ComplementItemID = nil
	groupItemModel.ComplementItem = nil

	groupItem := groupItemModel.ToDomain()
	groupItem.CalculateTotalPrice()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItem(ctx, groupItemModel); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	return nil
}

func (s *GroupItemService) UpdateGroupItemTotal(ctx context.Context, id string) error {
	groupItemModel, err := s.r.GetGroupByID(ctx, id, true)
	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()

	groupItem.CalculateTotalPrice()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItem(ctx, groupItemModel); err != nil {
		return err
	}

	return nil
}
