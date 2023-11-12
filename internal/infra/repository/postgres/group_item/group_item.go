package groupitemrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
)

type GroupItemRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewGroupItemRepositoryBun(db *bun.DB) *GroupItemRepositoryBun {
	return &GroupItemRepositoryBun{db: db}
}

func (r *GroupItemRepositoryBun) CreateGroupItem(ctx context.Context, p *itementity.GroupItem) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(p).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) UpdateGroupItem(ctx context.Context, p *itementity.GroupItem) error {
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
	_, err := r.db.NewDelete().Model(&itementity.GroupItem{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) GetGroupByID(ctx context.Context, id string) (*itementity.GroupItem, error) {
	item := &itementity.GroupItem{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *GroupItemRepositoryBun) GetAllPendingGroups(ctx context.Context) ([]itementity.GroupItem, error) {
	items := []itementity.GroupItem{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&items).Where("status = ?", "Pending").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *GroupItemRepositoryBun) GetGroupsByOrderId(ctx context.Context, id string) ([]itementity.GroupItem, error) {
	items := []itementity.GroupItem{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&items).Where("status = ?", "Pending").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return items, nil
}
