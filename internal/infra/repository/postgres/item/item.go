package itemrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type ItemRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewItemRepositoryBun(db *bun.DB) *ItemRepositoryBun {
	return &ItemRepositoryBun{db: db}
}

func (r *ItemRepositoryBun) AddItem(ctx context.Context, p *orderentity.Item) error {
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

func (r *ItemRepositoryBun) AddAdditionalItem(ctx context.Context, id uuid.UUID, productID uuid.UUID, additionalItem *orderentity.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	itemToAdditional := &orderentity.ItemToAdditional{
		ItemID:           id,
		AdditionalItemID: additionalItem.ID,
		ProductID:        productID,
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewDelete().Model(&orderentity.ItemToAdditional{}).Where("item_id = ? AND product_id = ?", id, productID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.NewInsert().Model(additionalItem).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.NewInsert().Model(itemToAdditional).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ItemRepositoryBun) UpdateItem(ctx context.Context, p *orderentity.Item) error {
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

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	// Apaga o item
	if _, err := tx.NewDelete().Model(&orderentity.Item{}).Where("id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	additionalItems := []orderentity.ItemToAdditional{}
	if err := r.db.NewSelect().Model(&additionalItems).Where("item_id = ?", id).Scan(ctx); err != nil {
		tx.Rollback()
		return err
	}

	additionalIds := []uuid.UUID{}
	for _, item := range additionalItems {
		additionalIds = append(additionalIds, item.AdditionalItemID)
	}

	// Apaga a relacao do item com additional items
	if _, err := tx.NewDelete().Model(&orderentity.ItemToAdditional{}).Where("item_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	// Apaga os additional items
	if len(additionalIds) > 0 {
		if _, err := tx.NewDelete().Model(&orderentity.Item{}).Where("id in (?)", bun.In(additionalIds)).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ItemRepositoryBun) DeleteAdditionalItem(ctx context.Context, idAdditional uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewDelete().Model(&orderentity.Item{}).Where("id = ?", idAdditional).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.NewDelete().Model(&orderentity.ItemToAdditional{}).Where("additional_item_id = ?", idAdditional).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ItemRepositoryBun) GetItemById(ctx context.Context, id string) (*orderentity.Item, error) {
	item := &orderentity.Item{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(item).Where("item.id = ?", id).Relation("AdditionalItems").Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepositoryBun) GetAllItems(ctx context.Context) ([]orderentity.Item, error) {
	items := []orderentity.Item{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&items).Relation("AdditionalItems").Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
