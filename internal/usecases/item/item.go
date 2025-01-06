package itemusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
)

var (
	ErrCategoryNotFound           = errors.New("category not found")
	ErrSizeNotFound               = errors.New("size not found")
	ErrSizeMustBeTheSame          = errors.New("size must be the same")
	ErrGroupItemNotBelongsToOrder = errors.New("group item not belongs to order")
	ErrGroupNotStaging            = errors.New("group not staging")
	ErrItemNotStagingAndPending   = errors.New("item not staging or pending")
)

type Service struct {
	ri  orderentity.ItemRepository
	rgi orderentity.GroupItemRepository
	ro  orderentity.OrderRepository
	rp  productentity.ProductRepository
	rq  productentity.QuantityRepository
	rc  productentity.CategoryRepository
}

func NewService(ri orderentity.ItemRepository) *Service {
	return &Service{ri: ri}
}

func (s *Service) AddDependencies(rgi orderentity.GroupItemRepository, ro orderentity.OrderRepository, rp productentity.ProductRepository, rq productentity.QuantityRepository, rc productentity.CategoryRepository) {
	s.rgi = rgi
	s.ro = ro
	s.rp = rp
	s.rq = rq
	s.rc = rc
}

func (s *Service) AddItemOrder(ctx context.Context, dto *itemdto.AddItemOrderInput) (ids *itemdto.ItemIDAndGroupItemDTO, err error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.ro.GetOrderById(ctx, dto.OrderID.String()); err != nil {
		return nil, err
	}

	product, err := s.rp.GetProductById(ctx, dto.ProductID.String())

	if err != nil {
		return nil, errors.New("product not found: " + err.Error())
	}

	if product.Category == nil {
		return nil, ErrCategoryNotFound
	}

	if product.Size == nil {
		return nil, ErrSizeNotFound
	}

	if dto.GroupItemID == nil {
		groupItem, err := s.newGroupItem(ctx, dto.OrderID, product)

		if err != nil {
			return nil, errors.New("new group item error: " + err.Error())
		}

		dto.GroupItemID = &groupItem.ID
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, dto.GroupItemID.String(), true)

	if err != nil {
		return nil, errors.New("group item not found: " + err.Error())
	}

	if groupItem.OrderID != dto.OrderID {
		return nil, ErrGroupItemNotBelongsToOrder
	}

	if ok, err := groupItem.CanAddItems(); !ok {
		return nil, err
	}

	quantity, err := s.rq.GetQuantityById(ctx, dto.QuantityID.String())

	if err != nil {
		return nil, errors.New("quantity not found: " + err.Error())
	}

	item, err := dto.ToModel(product, groupItem, quantity)

	if err != nil {
		return nil, err
	}

	if err = s.ri.AddItem(ctx, item); err != nil {
		return nil, errors.New("add item error: " + err.Error())
	}

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return nil, errors.New("update group item error: " + err.Error())
	}

	if groupItem.ComplementItemID != nil {
		groupItem.ComplementItem.Quantity += item.Quantity
		if err = s.ri.UpdateItem(ctx, groupItem.ComplementItem); err != nil {
			return nil, errors.New("update complement item error: " + err.Error())
		}
	}

	return itemdto.NewOutput(item.ID, groupItem.ID), nil
}

func (s *Service) DeleteItemOrder(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = s.ri.DeleteItem(ctx, dto.ID.String()); err != nil {
		return errors.New("delete item error: " + err.Error())
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)

	if err != nil {
		return errors.New("group item not found: " + err.Error())
	}

	if len(groupItem.Items) == 0 {
		var complementItemID *string
		if groupItem.ComplementItemID != nil {
			complementItemID = new(string)
			*complementItemID = groupItem.ComplementItemID.String()
		}

		if err = s.rgi.DeleteGroupItem(ctx, groupItem.ID.String(), complementItemID); err != nil {
			return err
		}

		return nil
	}

	if groupItem.ComplementItemID != nil {
		groupItem.ComplementItem.Quantity -= item.Quantity
		if err = s.ri.UpdateItem(ctx, groupItem.ComplementItem); err != nil {
			return errors.New("update complement item error: " + err.Error())
		}
	}

	return nil
}

func (s *Service) AddAdditionalItemOrder(ctx context.Context, dto *entitydto.IdRequest, dtoAdditional *itemdto.AddAdditionalItemOrderInput) (id uuid.UUID, err error) {
	productID, quantityID, err := dtoAdditional.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return uuid.Nil, errors.New("item not found: " + err.Error())
	}

	productAdditional, err := s.rp.GetProductById(ctx, productID.String())

	if err != nil {
		return uuid.Nil, errors.New("product not found: " + err.Error())
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)
	if err != nil {
		return uuid.Nil, errors.New("group item not found: " + err.Error())
	}

	if ok, err := groupItem.CanAddItems(); !ok {
		return uuid.Nil, err
	}

	found := false
	for _, additionalCategory := range groupItem.Category.AdditionalCategories {
		if additionalCategory.ID == productAdditional.CategoryID {
			found = true
			break
		}
	}

	if !found {
		return uuid.Nil, errors.New("additional category does not belong to this category")
	}

	quantity, err := s.rq.GetQuantityById(ctx, quantityID.String())

	if err != nil {
		return uuid.Nil, errors.New("quantity not found: " + err.Error())
	}

	if productAdditional.CategoryID != quantity.CategoryID {
		return uuid.Nil, errors.New("product category and quantity not match")
	}

	itemAdditional := orderentity.NewItem(productAdditional.Name, productAdditional.Price, quantity.Quantity, item.Size, productAdditional.ID, productAdditional.CategoryID)

	if err = s.ri.AddAdditionalItem(ctx, item.ID, productAdditional.ID, itemAdditional); err != nil {
		return uuid.Nil, errors.New("add additional item error: " + err.Error())
	}

	groupItem, err = s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)

	if err != nil {
		return uuid.Nil, errors.New("group item not found: " + err.Error())
	}

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return uuid.Nil, errors.New("update group item error: " + err.Error())
	}

	return itemAdditional.ID, nil
}

func (s *Service) DeleteAdditionalItemOrder(ctx context.Context, dtoAdditional *entitydto.IdRequest) (err error) {
	if err = s.ri.DeleteAdditionalItem(ctx, dtoAdditional.ID); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddRemovedItem(ctx context.Context, dtoID *entitydto.IdRequest, dto *itemdto.RemovedItemDTO) (err error) {
	name, err := dto.ToModel()
	if err != nil {
		return err
	}

	item, err := s.ri.GetItemById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

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

	return s.ri.UpdateItem(ctx, item)
}

func (s *Service) RemoveRemovedItem(ctx context.Context, dtoID *entitydto.IdRequest, dto *itemdto.RemovedItemDTO) (err error) {
	name, err := dto.ToModel()
	if err != nil {
		return err
	}

	item, err := s.ri.GetItemById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	item.RemoveRemovedItem(*name)

	return s.ri.UpdateItem(ctx, item)
}

func (s *Service) newGroupItem(ctx context.Context, orderID uuid.UUID, product *productentity.Product) (groupItem *orderentity.GroupItem, err error) {
	groupCommonAttributes := orderentity.GroupCommonAttributes{
		OrderID: orderID,
		GroupDetails: orderentity.GroupDetails{
			CategoryID:     product.CategoryID,
			Size:           product.Size.Name,
			NeedPrint:      product.Category.NeedPrint,
			UseProcessRule: product.Category.UseProcessRule,
		},
	}

	groupItem = orderentity.NewGroupItem(groupCommonAttributes)
	err = s.rgi.CreateGroupItem(ctx, groupItem)
	return
}
