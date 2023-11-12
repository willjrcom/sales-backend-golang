package itemrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
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
	_, err := r.db.NewInsert().Model(p).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepositoryBun) UpdateItem(ctx context.Context, p *itementity.Item) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepositoryBun) DeleteItem(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&itementity.Item{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepositoryBun) GetItemById(ctx context.Context, id string) (*itementity.Item, error) {
	item := &itementity.Item{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(item).Where("item.id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepositoryBun) GetAllItems(ctx context.Context) ([]itementity.Item, error) {
	items := []itementity.Item{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&items).Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return items, nil
}
