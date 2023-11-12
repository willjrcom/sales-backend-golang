package itementity

import (
	"context"
)

type ItemRepository interface {
	AddItem(ctx context.Context, item *Item) error
	DeleteItem(ctx context.Context, id string) error
	UpdateItem(ctx context.Context, item *Item) error
}

type GroupItemRepository interface {
	CreateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	GetGroupByID(ctx context.Context, id string) (*GroupItem, error)
	DeleteGroupItem(ctx context.Context, id string) error
	GetGroupsByOrderId(ctx context.Context, id string) ([]GroupItem, error)
	GetAllPendingGroups(ctx context.Context) ([]GroupItem, error)
}
