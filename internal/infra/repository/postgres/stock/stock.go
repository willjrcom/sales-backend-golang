package stockrepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type StockRepositoryBun struct {
	db *bun.DB
}

func NewStockRepositoryBun(db *bun.DB) model.StockRepository {
	return &StockRepositoryBun{db: db}
}

func (r *StockRepositoryBun) CreateStock(ctx context.Context, s *model.Stock) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockRepositoryBun) UpdateStock(ctx context.Context, s *model.Stock) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockRepositoryBun) GetStockByID(ctx context.Context, id string) (*model.Stock, error) {
	stock := &model.Stock{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(stock).Where("stock.id = ?", id).Relation("Product").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return stock, nil
}

func (r *StockRepositoryBun) GetStockByProductID(ctx context.Context, productID string) (*model.Stock, error) {
	stock := &model.Stock{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(stock).Where("stock.product_id = ?", productID).Relation("Product").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return stock, nil
}

func (r *StockRepositoryBun) GetAllStocks(ctx context.Context, page, perPage int) ([]model.Stock, int, error) {
	stocks := []model.Stock{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	query := tx.NewSelect().Model(&stocks).Relation("Product").Limit(perPage).Offset(page * perPage)

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return stocks, count, nil
}

func (r *StockRepositoryBun) GetActiveStocks(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&stocks).Where("stock.is_active = ?", true).Relation("Product").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepositoryBun) GetLowStockProducts(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&stocks).Where("stock.current_stock <= stock.min_stock AND stock.is_active = ?", true).Relation("Product").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepositoryBun) GetOutOfStockProducts(ctx context.Context) ([]model.Stock, error) {
	stocks := []model.Stock{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&stocks).Where("stock.current_stock <= 0 AND stock.is_active = ?", true).Relation("Product").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return stocks, nil
}
