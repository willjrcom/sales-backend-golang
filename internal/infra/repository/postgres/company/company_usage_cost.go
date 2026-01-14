package companyrepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyUsageCostRepository struct {
	db *bun.DB
}

func NewCompanyUsageCostRepository(db *bun.DB) *CompanyUsageCostRepository {
	return &CompanyUsageCostRepository{db: db}
}

func (r *CompanyUsageCostRepository) Create(ctx context.Context, cost *model.CompanyUsageCost) error {
	_, err := r.db.NewInsert().
		Model(cost).
		Exec(ctx)
	return err
}

func (r *CompanyUsageCostRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CompanyUsageCost, error) {
	cost := &model.CompanyUsageCost{}
	err := r.db.NewSelect().
		Model(cost).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return cost, nil
}

func (r *CompanyUsageCostRepository) GetMonthlyCosts(ctx context.Context, companyID uuid.UUID, month, year int) ([]*model.CompanyUsageCost, error) {
	var costs []*model.CompanyUsageCost
	err := r.db.NewSelect().
		Model(&costs).
		Where("company_id = ?", companyID).
		Where("billing_month = ?", month).
		Where("billing_year = ?", year).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return costs, nil
}

func (r *CompanyUsageCostRepository) GetCostsByType(ctx context.Context, companyID uuid.UUID, costType string, month, year int) ([]*model.CompanyUsageCost, error) {
	var costs []*model.CompanyUsageCost
	err := r.db.NewSelect().
		Model(&costs).
		Where("company_id = ?", companyID).
		Where("cost_type = ?", costType).
		Where("billing_month = ?", month).
		Where("billing_year = ?", year).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return costs, nil
}

func (r *CompanyUsageCostRepository) GetTotalByMonth(ctx context.Context, companyID uuid.UUID, month, year int) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := r.db.NewSelect().
		Model((*model.CompanyUsageCost)(nil)).
		ColumnExpr("COALESCE(SUM(amount), 0) as total").
		Where("company_id = ?", companyID).
		Where("billing_month = ?", month).
		Where("billing_year = ?", year).
		Where("deleted_at IS NULL").
		Scan(ctx, &total)
	if err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (r *CompanyUsageCostRepository) GetByReferenceID(ctx context.Context, referenceID uuid.UUID) (*model.CompanyUsageCost, error) {
	cost := &model.CompanyUsageCost{}
	err := r.db.NewSelect().
		Model(cost).
		Where("reference_id = ?", referenceID).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return cost, nil
}
