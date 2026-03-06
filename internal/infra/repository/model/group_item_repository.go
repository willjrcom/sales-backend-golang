package model

import (
	"context"

	"github.com/uptrace/bun"
)

type GroupItemRepository interface {
	CreateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	UpdateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	UpdateGroupItemWithTx(ctx context.Context, tx *bun.Tx, groupitem *GroupItem) (err error)
	GetGroupByID(ctx context.Context, id string, withRelation bool) (*GroupItem, error)
	GetGroupByIDWithTx(ctx context.Context, tx *bun.Tx, id string, withRelation bool) (*GroupItem, error)
	DeleteGroupItem(ctx context.Context, id string, complementItemID *string) error
	DeleteGroupItemWithTx(ctx context.Context, tx *bun.Tx, id string, complementItemID *string) error
	GetGroupItemsByOrderIDAndStatus(ctx context.Context, id string, status string) ([]GroupItem, error)
	GetGroupItemsByStatus(ctx context.Context, status string) ([]GroupItem, error)
	UpsertGroupItemSnapshot(ctx context.Context, snapshot *OrderGroupItemSnapshot) error
}
