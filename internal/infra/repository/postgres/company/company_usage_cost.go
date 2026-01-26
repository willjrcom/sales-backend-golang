package companyrepositorybun

import (
	"context"

	"github.com/google/uuid"
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

func (r *CompanyUsageCostRepository) GetPendingCosts(ctx context.Context, companyID uuid.UUID) ([]*model.CompanyUsageCost, error) {
	costs := []*model.CompanyUsageCost{}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	err = tx.NewSelect().
		Model(&costs).
		Where("company_id = ?", companyID).
		Where("status = ?", "PENDING").
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("payment_id IS NULL").
				WhereOr("payment_id IN (SELECT id FROM company_payments WHERE status IN ('refused', 'cancelled', 'rejected'))")
		}).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return costs, nil
}

func (r *CompanyUsageCostRepository) UpdateCostsPaymentID(ctx context.Context, costIDs []uuid.UUID, paymentID uuid.UUID) error {
	if len(costIDs) == 0 {
		return nil
	}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	_, err = tx.NewUpdate().
		Model((*model.CompanyUsageCost)(nil)).
		Set("payment_id = ?", paymentID).
		Where("id IN (?)", bun.In(costIDs)).
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
		Where("month = ?", month).
		Where("year = ?", year).
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

func (r *CompanyUsageCostRepository) Update(ctx context.Context, cost *model.CompanyUsageCost) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().
		Model(cost).
		WherePK().
		Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CompanyUsageCostRepository) GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*model.CompanyUsageCost, error) {
	var costs []*model.CompanyUsageCost
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(&costs).
		Where("payment_id = ?", paymentID).
		Where("deleted_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return costs, nil
}
