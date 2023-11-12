package groupitemrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
)

type GroupItemRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewGroupItemRepositoryBun(db *bun.DB) *GroupItemRepositoryBun {
	return &GroupItemRepositoryBun{db: db}
}

func (r *GroupItemRepositoryBun) CreateGroupItem(ctx context.Context, p *groupitementity.GroupItem) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(p).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) UpdateGroupItem(ctx context.Context, p *groupitementity.GroupItem) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) DeleteGroupItem(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&groupitementity.GroupItem{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) GetGroupByID(ctx context.Context, id string) (*groupitementity.GroupItem, error) {
	item := &groupitementity.GroupItem{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *GroupItemRepositoryBun) GetAllPendingGroups(ctx context.Context) ([]groupitementity.GroupItem, error) {
	items := []groupitementity.GroupItem{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&items).Where("status = ?", "Pending").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *GroupItemRepositoryBun) GetGroupsByOrderId(ctx context.Context, id string) ([]groupitementity.GroupItem, error) {
	items := []groupitementity.GroupItem{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&items).Where("status = ?", "Pending").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return items, nil
}
