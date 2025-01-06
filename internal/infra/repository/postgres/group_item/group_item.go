package groupitemrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type GroupItemRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewGroupItemRepositoryBun(db *bun.DB) *GroupItemRepositoryBun {
	return &GroupItemRepositoryBun{db: db}
}

func (r *GroupItemRepositoryBun) CreateGroupItem(ctx context.Context, p *model.GroupItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(p).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) UpdateGroupItem(ctx context.Context, p *model.GroupItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) DeleteGroupItem(ctx context.Context, id string, complementItemID *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewDelete().Model(&model.GroupItem{}).Where("id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if complementItemID != nil {
		if _, err = tx.NewDelete().Model(&model.Item{}).Where("id = ?", complementItemID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err = tx.NewDelete().Model(&model.Item{}).Where("group_item_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) GetGroupByID(ctx context.Context, id string, withRelation bool) (*model.GroupItem, error) {
	item := &model.GroupItem{}
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	query := r.db.NewSelect().Model(item).Where("group_item.id = ?", id).Relation("Category").Relation("ComplementItem")

	if withRelation {
		query.Relation("Items.AdditionalItems").Relation("Category.ComplementCategories").Relation("Category.AdditionalCategories")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *GroupItemRepositoryBun) GetGroupsByStatus(ctx context.Context, status model.StatusGroupItem) ([]model.GroupItem, error) {
	items := []model.GroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&items).Where("group_item.status = ?", status).Relation("Items.AdditionalItems").Relation("Category").Relation("ComplementItem").Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *GroupItemRepositoryBun) GetGroupsByOrderIDAndStatus(ctx context.Context, id string, status model.StatusGroupItem) ([]model.GroupItem, error) {
	items := []model.GroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&items).Where("group_item.order_id = ? AND group_item.status = ?", id, status).Relation("Items.AdditionalItems").Relation("ComplementItem").Relation("Category").Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
