package stockrepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type StockBatchRepositoryBun struct {
	db *bun.DB
}

func NewStockBatchRepositoryBun(db *bun.DB) model.StockBatchRepository {
	return &StockBatchRepositoryBun{db: db}
}

func (r *StockBatchRepositoryBun) CreateBatch(ctx context.Context, db bun.IDB, b *model.StockBatch) error {
	if _, err := db.NewInsert().Model(b).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *StockBatchRepositoryBun) UpdateBatch(ctx context.Context, db bun.IDB, b *model.StockBatch) error {
	if _, err := db.NewUpdate().Model(b).Where("stock_batch.id = ?", b.ID).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *StockBatchRepositoryBun) GetBatchByID(ctx context.Context, db bun.IDB, id string) (*model.StockBatch, error) {
	batch := &model.StockBatch{}

	if err := db.NewSelect().Model(batch).Where("stock_batch.id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return batch, nil
}

func (r *StockBatchRepositoryBun) GetBatchesByStockID(ctx context.Context, stockID string) ([]model.StockBatch, error) {
	batches := []model.StockBatch{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	// Fix #22: filtrar current_quantity > 0 para evitar alertas em lotes zerados.
	// Não filtramos expires_at aqui (diff. de GetActiveBatchesByStockID) para que
	// lotes VENCIDOS com estoque restante sejam visíveis ao CheckExpirations.
	if err := tx.NewSelect().Model(&batches).
		Where("stock_batch.stock_id = ?", stockID).
		Where("stock_batch.current_quantity > 0").
		Order("created_at ASC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *StockBatchRepositoryBun) GetActiveBatchesByStockID(ctx context.Context, stockID string) ([]model.StockBatch, error) {
	batches := []model.StockBatch{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&batches).
		Where("stock_batch.stock_id = ?", stockID).
		Where("stock_batch.current_quantity > 0").
		Where("stock_batch.expires_at IS NULL OR stock_batch.expires_at > NOW()").
		Order("created_at ASC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *StockBatchRepositoryBun) GetActiveBatchesByStockIDForUpdate(ctx context.Context, db bun.IDB, stockID string) ([]model.StockBatch, error) {
	batches := []model.StockBatch{}

	if err := db.NewSelect().Model(&batches).
		Where("stock_batch.stock_id = ?", stockID).
		Where("stock_batch.current_quantity > 0").
		Where("stock_batch.expires_at IS NULL OR stock_batch.expires_at > NOW()").
		Order("created_at ASC").
		For("UPDATE").
		Scan(ctx); err != nil {
		return nil, err
	}
	return batches, nil
}
