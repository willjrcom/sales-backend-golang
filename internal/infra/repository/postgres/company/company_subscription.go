package companyrepositorybun

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanySubscriptionRepositoryBun struct {
	db *bun.DB
}

func NewCompanySubscriptionRepositoryBun(db *bun.DB) model.CompanySubscriptionRepository {
	return &CompanySubscriptionRepositoryBun{db: db}
}

func (r *CompanySubscriptionRepositoryBun) CreateSubscription(ctx context.Context, subscription *model.CompanySubscription) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(subscription).Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CompanySubscriptionRepositoryBun) UpdateSubscription(ctx context.Context, subscription *model.CompanySubscription) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().
		Model(subscription).
		Where("id = ?", subscription.ID).
		Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CompanySubscriptionRepositoryBun) MarkActiveSubscriptionAsCanceled(ctx context.Context, companyID uuid.UUID) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	result, err := tx.NewUpdate().
		Model((*model.CompanySubscription)(nil)).
		Set("is_canceled = ?", true).
		Where("company_id = ?", companyID).
		Where("is_active = ?", true).
		Where("start_date <= NOW()"). // Must have already started (not upcoming)
		Where("end_date > NOW()").    // Must not have expired yet
		Exec(ctx)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no active subscription found to cancel")
	}

	return tx.Commit()
}

func (r *CompanySubscriptionRepositoryBun) GetActiveSubscription(ctx context.Context, companyID uuid.UUID) (*model.CompanySubscription, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	// Get Active
	active := &model.CompanySubscription{}
	if err := tx.NewSelect().
		Model(active).
		Where("company_id = ?", companyID).
		Where("is_active = ?", true).
		Where("end_date > ?", time.Now().UTC()).
		Order("end_date DESC").
		Limit(1).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			active = nil
		} else {
			return nil, err
		}
	}

	return active, nil
}

func (r *CompanySubscriptionRepositoryBun) UpdateCompanyPlans(ctx context.Context) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewRaw(`
		UPDATE company_subscriptions 
		SET is_active = false 
		WHERE is_active = true 
		AND end_date < NOW()
	`).Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
