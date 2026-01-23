package stockrepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type StockAlertRepositoryBun struct {
	db *bun.DB
}

func NewStockAlertRepositoryBun(db *bun.DB) model.StockAlertRepository {
	return &StockAlertRepositoryBun{db: db}
}

func (r *StockAlertRepositoryBun) CreateAlert(ctx context.Context, a *model.StockAlert) error {
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(a).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockAlertRepositoryBun) UpdateAlert(ctx context.Context, a *model.StockAlert) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(a).Where("id = ?", a.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockAlertRepositoryBun) GetAlertsByStockID(ctx context.Context, stockID string) ([]model.StockAlert, error) {
	alerts := []model.StockAlert{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&alerts).Where("stock_alert.stock_id = ?", stockID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *StockAlertRepositoryBun) GetActiveAlerts(ctx context.Context) ([]model.StockAlert, error) {
	alerts := []model.StockAlert{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&alerts).Where("stock_alert.is_resolved = ?", false).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *StockAlertRepositoryBun) GetResolvedAlerts(ctx context.Context) ([]model.StockAlert, error) {
	alerts := []model.StockAlert{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&alerts).Where("stock_alert.is_resolved = ?", true).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *StockAlertRepositoryBun) ResolveAlert(ctx context.Context, alertID string, resolvedBy string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	_, err = tx.NewUpdate().Model((*model.StockAlert)(nil)).
		Set("is_resolved = ?", true).
		Set("resolved_by = ?", resolvedBy).
		Set("resolved_at = NOW() ").
		Where("id = ?", alertID).
		Exec(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *StockAlertRepositoryBun) GetAlertByID(ctx context.Context, alertID string) (*model.StockAlert, error) {
	alert := &model.StockAlert{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(alert).
		Where("id = ?", alertID).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return alert, nil
}

func (r *StockAlertRepositoryBun) GetAllAlerts(ctx context.Context) ([]model.StockAlert, error) {
	var alerts []model.StockAlert

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(&alerts).
		Order("created_at DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return alerts, nil
}

func (r *StockAlertRepositoryBun) DeleteAlert(ctx context.Context, alertID string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().
		Model((*model.StockAlert)(nil)).
		Where("id = ?", alertID).
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
