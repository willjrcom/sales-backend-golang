package companyrepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyUsageCostRepository struct {
	db *bun.DB
}

func NewCompanyUsageCostRepository(db *bun.DB) *CompanyUsageCostRepository {
	return &CompanyUsageCostRepository{db: db}
}

func (r *CompanyUsageCostRepository) Create(ctx context.Context, cost *model.CompanyUsageCost) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	_, err = tx.NewInsert().
		Model(cost).
		Exec(ctx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *CompanyUsageCostRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CompanyUsageCost, error) {
	cost := &model.CompanyUsageCost{}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(cost).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return cost, nil
}

func (r *CompanyUsageCostRepository) GetMonthlyCosts(ctx context.Context, companyID uuid.UUID, month, year int) ([]*model.CompanyUsageCost, error) {
	var costs []*model.CompanyUsageCost
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(&costs).
		Where("company_id = ?", companyID).
		Where("billing_month = ?", month).
		Where("billing_year = ?", year).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return costs, nil
}

func (r *CompanyUsageCostRepository) GetCostsByType(ctx context.Context, companyID uuid.UUID, costType string, month, year int) ([]*model.CompanyUsageCost, error) {
	var costs []*model.CompanyUsageCost
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(&costs).
		Where("company_id = ?", companyID).
		Where("cost_type = ?", costType).
		Where("billing_month = ?", month).
		Where("billing_year = ?", year).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return costs, nil
}

func (r *CompanyUsageCostRepository) GetTotalByMonth(ctx context.Context, companyID uuid.UUID, month, year int) (decimal.Decimal, error) {
	var total decimal.Decimal
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return decimal.Zero, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model((*model.CompanyUsageCost)(nil)).
		ColumnExpr("COALESCE(SUM(amount), 0) as total").
		Where("company_id = ?", companyID).
		Where("billing_month = ?", month).
		Where("billing_year = ?", year).
		Where("deleted_at IS NULL").
		Scan(ctx, &total); err != nil {
		return decimal.Zero, err
	}

	if err := tx.Commit(); err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (r *CompanyUsageCostRepository) GetByReferenceID(ctx context.Context, referenceID uuid.UUID) (*model.CompanyUsageCost, error) {
	cost := &model.CompanyUsageCost{}
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(cost).
		Where("reference_id = ?", referenceID).
		Where("deleted_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return cost, nil
}
