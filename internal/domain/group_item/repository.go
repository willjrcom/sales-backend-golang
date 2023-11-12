package groupitementity

import (
	"context"
)

type GroupItemRepository interface {
	CreateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	GetGroupByID(ctx context.Context, id string) (*GroupItem, error)
	DeleteGroupItem(ctx context.Context, id string) error
	GetGroupsByOrderId(ctx context.Context, id string) ([]GroupItem, error)
	GetAllPendingGroups(ctx context.Context) ([]GroupItem, error)
}
