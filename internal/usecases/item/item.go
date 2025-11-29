package itemusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
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
	ri  model.ItemRepository
	rgi model.GroupItemRepository
	ro  model.OrderRepository
	rp  model.ProductRepository
	rq  model.QuantityRepository
	rc  model.CategoryRepository
	so  *orderusecases.OrderService
	sgi *orderusecases.GroupItemService
}

func NewService(ri model.ItemRepository) *Service {
	return &Service{ri: ri}
}

func (s *Service) AddDependencies(rgi model.GroupItemRepository, ro model.OrderRepository, rp model.ProductRepository, rq model.QuantityRepository, rc model.CategoryRepository, so *orderusecases.OrderService, sgi *orderusecases.GroupItemService) {
	s.rgi = rgi
	s.ro = ro
	s.rp = rp
	s.rq = rq
	s.rc = rc
	s.so = so
	s.sgi = sgi
}

func (s *Service) AddItemOrder(ctx context.Context, dto *itemdto.OrderItemCreateDTO) (ids *itemdto.ItemIDDTO, err error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.ro.GetOrderById(ctx, dto.OrderID.String()); err != nil {
		return nil, err
	}

	productModel, err := s.rp.GetProductById(ctx, dto.ProductID.String())
	if err != nil {
		return nil, errors.New("product not found: " + err.Error())
	}

	product := productModel.ToDomain()

	if product.Category == nil {
		return nil, ErrCategoryNotFound
	}

	if product.Size == nil {
		return nil, ErrSizeNotFound
	}

	// If group item id is not provided, create a new group item
	if dto.GroupItemID == nil {
		groupItem, err := s.newGroupItem(ctx, dto.OrderID, product)

		if err != nil {
			return nil, errors.New("new group item error: " + err.Error())
		}

		dto.GroupItemID = &groupItem.ID
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, dto.GroupItemID.String(), true)

	if err != nil {
		return nil, errors.New("group item not found: " + err.Error())
	}

	groupItem := groupItemModel.ToDomain()

	// Check if group item belongs to order
	if groupItem.OrderID != dto.OrderID {
		return nil, ErrGroupItemNotBelongsToOrder
	}

	if ok, err := groupItem.CanAddItems(); !ok {
		return nil, err
	}

	quantityModel, err := s.rq.GetQuantityById(ctx, dto.QuantityID.String())

	if err != nil {
		return nil, errors.New("quantity not found: " + err.Error())
	}

	quantity := quantityModel.ToDomain()
	item, err := dto.ToDomain(product, groupItem, quantity)

	if err != nil {
		return nil, err
	}

	itemModel := &model.Item{}
	itemModel.FromDomain(item)

	if err = s.ri.AddItem(ctx, itemModel); err != nil {
		return nil, errors.New("add item error: " + err.Error())
	}

	GroupItemModel := &model.GroupItem{}
	GroupItemModel.FromDomain(groupItem)

	if err = s.rgi.UpdateGroupItem(ctx, groupItemModel); err != nil {
		return nil, errors.New("update group item error: " + err.Error())
	}

	// Update complement item
	if groupItem.ComplementItemID != nil {
		groupItem.ComplementItem.Quantity += item.Quantity

		complementItemModel := &model.Item{}
		complementItemModel.FromDomain(groupItem.ComplementItem)
		if err = s.ri.UpdateItem(ctx, complementItemModel); err != nil {
			return nil, errors.New("update complement item error: " + err.Error())
		}
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, dto.GroupItemID.String()); err != nil {
		return nil, err
	}

	if err := s.so.UpdateOrderTotal(ctx, dto.OrderID.String()); err != nil {
		return nil, err
	}

	return itemdto.FromDomain(item.ID, groupItem.ID), nil
}

func (s *Service) DeleteItemOrder(ctx context.Context, dto *entitydto.IDRequest) (err error) {
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

	// Update complement item quantity
	if groupItem.ComplementItemID != nil && len(groupItem.Items) != 0 {
		groupItem.ComplementItem.Quantity -= item.Quantity
		if err = s.ri.UpdateItem(ctx, groupItem.ComplementItem); err != nil {
			return errors.New("update complement item error: " + err.Error())
		}
	}

	// Delete group item if there are no items
	if len(groupItem.Items) == 0 {
		var complementItemID *string
		if groupItem.ComplementItemID != nil {
			complementItemID = new(string)
			*complementItemID = groupItem.ComplementItemID.String()
		}

		if err = s.rgi.DeleteGroupItem(ctx, groupItem.ID.String(), complementItemID); err != nil {
			return err
		}

		if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
			return err
		}

		return nil
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, groupItem.ID.String()); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddAdditionalItemOrder(ctx context.Context, dto *entitydto.IDRequest, dtoAdditional *itemdto.OrderAdditionalItemCreateDTO) (id uuid.UUID, err error) {
	productID, quantityID, flavor, err := dtoAdditional.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return uuid.Nil, errors.New("item not found: " + err.Error())
	}

	productAdditionalModel, err := s.rp.GetProductById(ctx, productID.String())

	if err != nil {
		return uuid.Nil, errors.New("product not found: " + err.Error())
	}
	productAdditional := productAdditionalModel.ToDomain()

	groupItemModel, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)
	if err != nil {
		return uuid.Nil, errors.New("group item not found: " + err.Error())
	}

	groupItem := groupItemModel.ToDomain()

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

	quantityModel, err := s.rq.GetQuantityById(ctx, quantityID.String())

	if err != nil {
		return uuid.Nil, errors.New("quantity not found: " + err.Error())
	}
	quantity := quantityModel.ToDomain()

	if productAdditional.CategoryID != quantity.CategoryID {
		return uuid.Nil, errors.New("product category and quantity not match")
	}

	normalizedFlavor, err := itemdto.NormalizeFlavor(flavor, productAdditional.Flavors)
	if err != nil {
		return uuid.Nil, err
	}

	additionalItem := orderentity.NewItem(productAdditional.Name, productAdditional.Price, quantity.Quantity, item.Size, productAdditional.ID, productAdditional.CategoryID, normalizedFlavor)
	additionalItem.IsAdditional = true
	additionalItem.GroupItemID = groupItem.ID

	additionalItemModel := &model.Item{}
	additionalItemModel.FromDomain(additionalItem)

	if err = s.ri.AddAdditionalItem(ctx, item.ID, productAdditional.ID, additionalItemModel); err != nil {
		return uuid.Nil, errors.New("add additional item error: " + err.Error())
	}

	if err := s.UpdateItemTotal(ctx, item.ID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, groupItem.ID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return uuid.Nil, err
	}

	return additionalItem.ID, nil
}

func (s *Service) DeleteAdditionalItemOrder(ctx context.Context, dtoAdditional *entitydto.IDRequest) (err error) {
	additionalItem, err := s.ri.GetItemById(ctx, dtoAdditional.ID.String())
	if err != nil {
		return errors.New("item not found: " + err.Error())
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, additionalItem.GroupItemID.String(), true)
	if err != nil {
		return errors.New("group item not found: " + err.Error())
	}

	item, err := s.ri.GetItemByAdditionalItemID(ctx, additionalItem.ID)
	if err != nil {
		return err
	}

	if err = s.ri.DeleteAdditionalItem(ctx, dtoAdditional.ID); err != nil {
		return err
	}

	if err := s.UpdateItemTotal(ctx, item.ID.String()); err != nil {
		return err
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, groupItemModel.ID.String()); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddRemovedItem(ctx context.Context, dtoID *entitydto.IDRequest, dto *itemdto.RemovedItemDTO) (err error) {
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

func (s *Service) RemoveRemovedItem(ctx context.Context, dtoID *entitydto.IDRequest, dto *itemdto.RemovedItemDTO) (err error) {
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

func (s *Service) newGroupItem(ctx context.Context, orderID uuid.UUID, product *productentity.Product) (groupItem *orderentity.GroupItem, err error) {
	groupCommonAttributes := orderentity.GroupCommonAttributes{
		OrderID: orderID,
		GroupDetails: orderentity.GroupDetails{
			CategoryID:     product.CategoryID,
			Size:           product.Size.Name,
			NeedPrint:      product.Category.NeedPrint,
			PrinterName:    product.Category.PrinterName,
			UseProcessRule: product.Category.UseProcessRule,
		},
	}

	groupItem = orderentity.NewGroupItem(groupCommonAttributes)

	groupItemModel := &model.GroupItem{}
	groupItemModel.FromDomain(groupItem)
	err = s.rgi.CreateGroupItem(ctx, groupItemModel)
	return
}

func (s *Service) UpdateItemTotal(ctx context.Context, id string) error {
	itemModel, err := s.ri.GetItemById(ctx, id)
	if err != nil {
		return err
	}

	item := itemModel.ToDomain()

	item.CalculateTotalPrice()

	itemModel.FromDomain(item)
	if err := s.ri.UpdateItem(ctx, itemModel); err != nil {
		return err
	}

	return nil
}
