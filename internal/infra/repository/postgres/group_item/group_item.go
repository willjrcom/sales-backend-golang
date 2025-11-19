package groupitemrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type GroupItemRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewGroupItemRepositoryBun(db *bun.DB) model.GroupItemRepository {
	return &GroupItemRepositoryBun{db: db}
}

func (r *GroupItemRepositoryBun) CreateGroupItem(ctx context.Context, p *model.GroupItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(p).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *GroupItemRepositoryBun) UpdateGroupItem(ctx context.Context, p *model.GroupItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *GroupItemRepositoryBun) DeleteGroupItem(ctx context.Context, id string, complementItemID *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
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

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	query := tx.NewSelect().Model(item).Where("group_item.id = ?", id).Relation("Category").Relation("ComplementItem")

	if withRelation {
		query.
			Relation("Items", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("is_additional = ?", false)
			}).
			Relation("Items.AdditionalItems").
			Relation("Category.ComplementCategories").
			Relation("Category.AdditionalCategories")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *GroupItemRepositoryBun) GetGroupItemsByStatus(ctx context.Context, status string) ([]model.GroupItem, error) {
	items := []model.GroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&items).
		Where("group_item.status = ?", status).
		Relation("Items", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("is_additional = ?", false)
		}).
		Relation("Items.AdditionalItems").
		Relation("Category").
		Relation("ComplementItem").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupItemRepositoryBun) GetGroupItemsByOrderIDAndStatus(ctx context.Context, id string, status string) ([]model.GroupItem, error) {
	items := []model.GroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&items).
		Where("group_item.order_id = ? AND group_item.status = ?", id, status).
		Relation("Items", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("is_additional = ?", false)
		}).
		Relation("Items.AdditionalItems").
		Relation("ComplementItem").
		Relation("Category").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return items, nil
}
