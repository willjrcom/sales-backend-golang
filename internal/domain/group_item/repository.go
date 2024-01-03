package groupitementity

import (
	"context"
)

type GroupItemRepository interface {
	CreateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	UpdateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	CalculateTotal(ctx context.Context, id string) (err error)
	GetGroupByID(ctx context.Context, id string, withRelation bool) (*GroupItem, error)
	DeleteGroupItem(ctx context.Context, id string) error
	GetGroupsByOrderID(ctx context.Context, id string) ([]GroupItem, error)
	GetAllGroupsByStatus(ctx context.Context, status StatusGroupItem) ([]GroupItem, error)
}
