package itementity

import (
	"context"
)

type ItemRepository interface {
	AddItemOrder(ctx context.Context, item *Item) error
	RemoveItemOrder(ctx context.Context, item *Item) error
	UpdateItemOrder(ctx context.Context, id string, item *Item) error
}

type GroupItemRepository interface {
	CreateGroupItem(ctx context.Context, item *GroupItem) (err error)
	GetGroupsByOrderId(ctx context.Context, id string) ([]GroupItem, error)
	GetGroupItemByOrderIdAndCategoryID(ctx context.Context, orderID, categoryID string) ([]GroupItem, error)
	GetGroupItemByID(ctx context.Context, id string) (*GroupItem, error)
}
