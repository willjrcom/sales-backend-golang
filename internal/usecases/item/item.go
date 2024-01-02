package itemusecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
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

func (s *Service) AddItemOrder(ctx context.Context, dto *itemdto.AddItemOrderInput) (id uuid.UUID, err error) {
	if _, err := s.ro.GetOrderById(ctx, dto.OrderID.String()); err != nil {
		return uuid.Nil, err
	}

	product, err := s.rp.GetProductById(ctx, dto.ProductID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if dto.GroupItemID == nil {
		groupItem, err := s.newGroupItem(ctx, dto.OrderID, product)

		if err != nil {
			return uuid.Nil, err
		}

		dto.GroupItemID = &groupItem.ID
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, dto.GroupItemID.String(), false)

	if err != nil {
		return uuid.Nil, err
	}

	quantity, err := s.rq.GetQuantityById(ctx, dto.QuantityID.String())

	if err != nil {
		return uuid.Nil, err
	}

	item, err := dto.ToModel(product, groupItem, quantity)

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.ri.AddItem(ctx, item); err != nil {
		return uuid.Nil, err
	}

	if err = s.rgi.CalculateTotal(ctx, groupItem.ID.String()); err != nil {
		return uuid.Nil, err
	}

	return item.ID, nil
}

func (s *Service) RemoveItemOrder(ctx context.Context, dto *entitydto.IdRequest) (err error) {
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

func (s *Service) AddAditionalItemOrder(ctx context.Context, dto *entitydto.IdRequest) (id uuid.UUID, err error) {
	return uuid.Nil, nil
}

func (s *Service) newGroupItem(ctx context.Context, orderID uuid.UUID, product *productentity.Product) (groupItem *groupitementity.GroupItem, err error) {
	groupCommonAttributes := groupitementity.GroupCommonAttributes{
		OrderID: orderID,
		GroupDetails: groupitementity.GroupDetails{
			CategoryID: product.CategoryID,
			Status:     groupitementity.StatusGroupStaging,
			Size:       product.Size.Name,
		},
	}

	groupItem = &groupitementity.GroupItem{
		Entity:                entity.NewEntity(),
		GroupCommonAttributes: groupCommonAttributes,
	}

	err = s.rgi.CreateGroupItem(ctx, groupItem)
	return
}
