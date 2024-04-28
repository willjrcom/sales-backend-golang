package itemusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
)

var (
	ErrCategoryNotFound         = errors.New("category not found")
	ErrSizeNotFound             = errors.New("size not found")
	ErrSizeMustBeTheSame        = errors.New("size must be the same")
	ErrGroupNotStaging          = errors.New("group not staging")
	ErrItemNotStagingAndPending = errors.New("item not staging or pending")
)

type Service struct {
	ri  itementity.ItemRepository
	rgi groupitementity.GroupItemRepository
	ro  orderentity.OrderRepository
	rp  productentity.ProductRepository
	rq  productentity.QuantityRepository
}

func NewService(ri itementity.ItemRepository, rgi groupitementity.GroupItemRepository, ro orderentity.OrderRepository, rp productentity.ProductRepository, rq productentity.QuantityRepository) *Service {
	return &Service{ri: ri, rgi: rgi, ro: ro, rp: rp, rq: rq}
}

func (s *Service) AddItemOrder(ctx context.Context, dto *itemdto.AddItemOrderInput) (ids *itemdto.ItemIDAndGroupItemOutput, err error) {
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

	if !groupItem.CanAddItems() {
		return nil, ErrGroupNotStaging
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

	groupItem.CalculateTotalPrice()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return nil, errors.New("update group item error: " + err.Error())
	}

	if groupItem.ComplementItemID != nil {
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

	groupItem.CalculateTotalPrice()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return errors.New("update group itemerror: " + err.Error())
	}

	if err = s.ri.UpdateItem(ctx, groupItem.ComplementItem); err != nil {
		return errors.New("update complement item error: " + err.Error())
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

	if !item.CanAddAdditionalItems() {
		return uuid.Nil, ErrItemNotStagingAndPending
	}

	productAdditional, err := s.rp.GetProductById(ctx, productID.String())

	if err != nil {
		return uuid.Nil, errors.New("product not found: " + err.Error())
	}

	groupItem, err := s.rgi.GetGroupByIDWithCategoryComplete(ctx, item.GroupItemID.String())
	if err != nil {
		return uuid.Nil, errors.New("group item not found: " + err.Error())
	}

	found := false
	for _, additionalCategory := range groupItem.Category.AdditionalCategories {
		if additionalCategory.ID == productAdditional.CategoryID {
			found = true
			break
		}
	}

	if !found {
		return uuid.Nil, errors.New("category product and additional not match")
	}

	quantity, err := s.rq.GetQuantityById(ctx, quantityID.String())

	if err != nil {
		return uuid.Nil, errors.New("quantity not found: " + err.Error())
	}

	if productAdditional.CategoryID != quantity.CategoryID {
		return uuid.Nil, errors.New("category product and quantity not match")
	}

	itemAdditional := itementity.NewItem(productAdditional.Name, productAdditional.Price, quantity.Quantity, item.Size, item.Status, productAdditional.ID)

	if err = s.ri.AddAdditionalItem(ctx, item.ID, itemAdditional); err != nil {
		return uuid.Nil, errors.New("add additional item error: " + err.Error())
	}

	groupItem, err = s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)

	if err != nil {
		return uuid.Nil, errors.New("group item not found: " + err.Error())
	}

	groupItem.CalculateTotalPrice()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return uuid.Nil, errors.New("update group item error: " + err.Error())
	}

	return itemAdditional.ID, nil
}

func (s *Service) DeleteAdditionalItemOrder(ctx context.Context, dto *entitydto.IdRequest, dtoAdditional *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = s.ri.DeleteAdditionalItem(ctx, item.ID, dtoAdditional.ID); err != nil {
		return err
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)

	if err != nil {
		return err
	}

	groupItem.CalculateTotalPrice()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	return nil
}

func (s *Service) newGroupItem(ctx context.Context, orderID uuid.UUID, product *productentity.Product) (groupItem *groupitementity.GroupItem, err error) {
	groupCommonAttributes := groupitementity.GroupCommonAttributes{
		OrderID: orderID,
		GroupDetails: groupitementity.GroupDetails{
			CategoryID: product.CategoryID,
			Size:       product.Size.Name,
			NeedPrint:  product.Category.NeedPrint,
		},
	}

	groupItem = groupitementity.NewGroupItem(groupCommonAttributes)
	err = s.rgi.CreateGroupItem(ctx, groupItem)
	return
}
