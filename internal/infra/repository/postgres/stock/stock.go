package stockrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type StockRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewStockRepositoryBun(db *bun.DB) model.StockRepository {
	return &StockRepositoryBun{db: db}
}

func (r *StockRepositoryBun) CreateStock(ctx context.Context, s *model.Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockRepositoryBun) UpdateStock(ctx context.Context, s *model.Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockRepositoryBun) GetStockByID(ctx context.Context, id string) (*model.Stock, error) {
	stock := &model.Stock{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(stock).Where("stock.id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return stock, nil
}

func (r *StockRepositoryBun) GetStockByProductID(ctx context.Context, productID string) (*model.Stock, error) {
	stock := &model.Stock{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(stock).Where("stock.product_id = ?", productID).Scan(ctx); err != nil {
		return nil, err
	}

	return stock, nil
}

func (r *StockRepositoryBun) GetAllStocks(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&stocks).Scan(ctx); err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *StockRepositoryBun) GetActiveStocks(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&stocks).Where("stock.is_active = ?", true).Scan(ctx); err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *StockRepositoryBun) GetLowStockProducts(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&stocks).Where("stock.current_stock <= stock.min_stock AND stock.is_active = ?", true).Scan(ctx); err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *StockRepositoryBun) GetOutOfStockProducts(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&stocks).Where("stock.current_stock <= 0 AND stock.is_active = ?", true).Scan(ctx); err != nil {
		return nil, err
	}

	return stocks, nil
}
