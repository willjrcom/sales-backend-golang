package groupitemrepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type GroupItemRepositoryLocal struct {}

func NewGroupItemRepositoryLocal() model.GroupItemRepository {
	return &GroupItemRepositoryLocal{}
}

func (r *GroupItemRepositoryLocal) CreateGroupItem(ctx context.Context, groupitem *model.GroupItem) error {
	return nil
}

func (r *GroupItemRepositoryLocal) UpdateGroupItem(ctx context.Context, groupitem *model.GroupItem) error {
	return nil
}

func (r *GroupItemRepositoryLocal) GetGroupByID(ctx context.Context, id string, withRelation bool) (*model.GroupItem, error) {
	return nil, nil
}

func (r *GroupItemRepositoryLocal) DeleteGroupItem(ctx context.Context, id string, complementItemID *string) error {
	return nil
}

func (r *GroupItemRepositoryLocal) GetGroupItemsByOrderIDAndStatus(ctx context.Context, id string, status string) ([]model.GroupItem, error) {
	return nil, nil
}

func (r *GroupItemRepositoryLocal) GetGroupItemsByStatus(ctx context.Context, status string) ([]model.GroupItem, error) {
	return nil, nil
}
