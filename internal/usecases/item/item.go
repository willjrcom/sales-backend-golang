package itemusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
)

var (
	ErrGroupItemNotStaging      = errors.New("group item not staging")
	ErrGroupItemCategoryInvalid = errors.New("group item category invalid")
)

type Service struct {
	ri  itementity.ItemRepository
	rgi itementity.GroupItemRepository
	ro  orderentity.OrderRepository
	rp  productentity.ProductRepository
}

func NewService(ri itementity.ItemRepository, rgi itementity.GroupItemRepository, ro orderentity.OrderRepository, rp productentity.ProductRepository) *Service {
	return &Service{ri: ri, rgi: rgi, ro: ro, rp: rp}
}

func (s *Service) AddItemOrder(ctx context.Context, idOrder, idProduct string, dto *itemdto.AddItemOrderInput) (id uuid.UUID, err error) {
	order, err := s.ro.GetOrderById(ctx, dto.OrderID.String())

	if err != nil {
		return uuid.Nil, err
	}

	product, err := s.rp.GetProductById(ctx, dto.ProductID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if dto.GroupItemID == nil {
		groupItem, err := s.newGroupItem(ctx, order, product)

		if err != nil {
			return uuid.Nil, err
		}

		dto.GroupItemID = &groupItem.ID
	}

	groupItem, err := s.rgi.GetGroupItemByID(ctx, dto.GroupItemID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if groupItem.Status != itementity.StatusItemStaging {
		return uuid.Nil, ErrGroupItemNotStaging
	}

	if groupItem.CategoryID != product.CategoryID {
		return uuid.Nil, ErrGroupItemCategoryInvalid
	}

	item, err := dto.ToModel(product)

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.ri.AddItemOrder(ctx, item); err != nil {
		return uuid.Nil, err
	}

	return item.ID, nil
}

func (s *Service) newGroupItem(ctx context.Context, order *orderentity.Order, product *productentity.Product) (groupItem *itementity.GroupItem, err error) {
	groupItem = &itementity.GroupItem{
		Entity:     entity.NewEntity(),
		OrderID:    order.ID,
		CategoryID: product.CategoryID,
		Status:     itementity.StatusItemStaging,
	}

	err = s.rgi.CreateGroupItem(ctx, groupItem)
	return
}
