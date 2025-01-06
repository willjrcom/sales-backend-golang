package model

import "context"

type GroupItemRepository interface {
	CreateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	UpdateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	GetGroupByID(ctx context.Context, id string, withRelation bool) (*GroupItem, error)
	DeleteGroupItem(ctx context.Context, id string, complementItemID *string) error
	GetGroupsByOrderIDAndStatus(ctx context.Context, id string, status string) ([]GroupItem, error)
	GetGroupsByStatus(ctx context.Context, status string) ([]GroupItem, error)
}
