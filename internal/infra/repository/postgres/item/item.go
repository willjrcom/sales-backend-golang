package itemrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
)

type ItemRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewItemRepositoryBun(db *bun.DB) *ItemRepositoryBun {
	return &ItemRepositoryBun{db: db}
}

func (r *ItemRepositoryBun) AddItem(ctx context.Context, p *itementity.Item) error {
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

func (r *ItemRepositoryBun) UpdateItem(ctx context.Context, p *itementity.Item) error {
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

func (r *ItemRepositoryBun) DeleteItem(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&itementity.Item{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ItemRepositoryBun) GetItemById(ctx context.Context, id string) (*itementity.Item, error) {
	item := &itementity.Item{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(item).Where("item.id = ?", id).Relation("AdditionalItem").Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepositoryBun) GetAllItems(ctx context.Context) ([]itementity.Item, error) {
	items := []itementity.Item{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&items).Relation("AdditionalItem").Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
