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
		return nil, err
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
			return nil, err
		}

		dto.GroupItemID = &groupItem.ID
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, dto.GroupItemID.String(), false)

	if err != nil {
		return nil, err
	}

	if !groupItem.CanAddItems() {
		return nil, ErrGroupNotStaging
	}

	quantity, err := s.rq.GetQuantityById(ctx, dto.QuantityID.String())

	if err != nil {
		return nil, err
	}

	item, err := dto.ToModel(product, groupItem, quantity)

	if err != nil {
		return nil, err
	}

	if err = s.ri.AddItem(ctx, item); err != nil {
		return nil, err
	}

	groupItem.CalculateTotalValues()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return nil, err
	}

	return &itemdto.ItemIDAndGroupItemOutput{
		GroupItemID: groupItem.ID,
		ItemID:      item.ID,
	}, nil
}

func (s *Service) DeleteItemOrder(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = s.ri.DeleteItem(ctx, dto.ID.String()); err != nil {
		return err
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)

	if err != nil {
		return err
	}

	if len(groupItem.Items) == 0 {
		if err = s.rgi.DeleteGroupItem(ctx, item.GroupItemID.String()); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) AddAdditionalItemOrder(ctx context.Context, dto *entitydto.IdRequest, dtoAdditional *entitydto.IdRequest) (id uuid.UUID, err error) {

	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if !item.CanAddAdditionalItems() {
		return uuid.Nil, ErrItemNotStagingAndPending
	}

	productAdditional, err := s.rp.GetProductById(ctx, dtoAdditional.ID.String())

	if err != nil {
		return uuid.Nil, err
	}

	itemAdditionalCommonAttributes := itementity.ItemCommonAttributes{
		Name:     productAdditional.Name,
		Status:   item.Status,
		Price:    productAdditional.Price * item.Quantity,
		Size:     item.Size,
		Quantity: item.Quantity,
	}

	itemAdditional := itementity.NewItem(itemAdditionalCommonAttributes)

	if err = s.ri.AddAdditionalItem(ctx, item.ID, itemAdditional); err != nil {
		return uuid.Nil, err
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)

	if err != nil {
		return uuid.Nil, err
	}

	groupItem.CalculateTotalValues()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return uuid.Nil, err
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

	groupItem.CalculateTotalValues()

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
