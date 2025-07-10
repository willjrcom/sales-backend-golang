package stockrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type StockMovementRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewStockMovementRepositoryBun(db *bun.DB) model.StockMovementRepository {
	return &StockMovementRepositoryBun{db: db}
}

func (r *StockMovementRepositoryBun) CreateMovement(ctx context.Context, m *model.StockMovement) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(m).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockMovementRepositoryBun) GetMovementsByStockID(ctx context.Context, stockID string) ([]model.StockMovement, error) {
	movements := []model.StockMovement{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&movements).Where("stock_movement.stock_id = ?", stockID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return movements, nil
}

func (r *StockMovementRepositoryBun) GetMovementsByProductID(ctx context.Context, productID string) ([]model.StockMovement, error) {
	movements := []model.StockMovement{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&movements).Where("stock_movement.product_id = ?", productID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return movements, nil
}

func (r *StockMovementRepositoryBun) GetMovementsByOrderID(ctx context.Context, orderID string) ([]model.StockMovement, error) {
	movements := []model.StockMovement{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&movements).Where("stock_movement.order_id = ?", orderID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return movements, nil
}

func (r *StockMovementRepositoryBun) GetAllMovements(ctx context.Context) ([]model.StockMovement, error) {
	movements := []model.StockMovement{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&movements).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return movements, nil
}

func (r *StockMovementRepositoryBun) GetMovementsByDateRange(ctx context.Context, start, end string) ([]model.StockMovement, error) {
	movements := []model.StockMovement{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&movements).Where("stock_movement.created_at BETWEEN ? AND ?", start, end).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return movements, nil
}
