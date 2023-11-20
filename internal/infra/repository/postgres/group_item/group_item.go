package groupitemrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/uptrace/bun"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
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
	p.CalculateTotalValues()

	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) CalculateTotal(ctx context.Context, id string) (err error) {
	groupItem, err := r.GetGroupByID(ctx, id, true)

	if err != nil {
		return err
	}

	groupItem.CalculateTotalValues()

	return r.UpdateGroupItem(ctx, groupItem)
}

func (r *GroupItemRepositoryBun) DeleteGroupItem(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	_, err = tx.NewDelete().Model(&groupitementity.GroupItem{}).Where("id = ?", id).Exec(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NewDelete().Model(&itementity.Item{}).Where("group_item_id = ?", id).Exec(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *GroupItemRepositoryBun) GetGroupByID(ctx context.Context, id string, withRelation bool) (*groupitementity.GroupItem, error) {
	item := &groupitementity.GroupItem{}

	r.mu.Lock()
	query := r.db.NewSelect().Model(item).Where("id = ?", id)

	if withRelation {
		query.Relation("Items")
	}

	err := query.Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *GroupItemRepositoryBun) GetAllPendingGroups(ctx context.Context) ([]groupitementity.GroupItem, error) {
	items := []groupitementity.GroupItem{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&items).Where("status = ?", "Pending").Relation("Items").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *GroupItemRepositoryBun) GetGroupsByOrderId(ctx context.Context, id string) ([]groupitementity.GroupItem, error) {
	items := []groupitementity.GroupItem{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&items).Where("status = ?", "Pending").Relation("Items").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return items, nil
}
