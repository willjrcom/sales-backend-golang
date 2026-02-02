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

func (r *CompanySubscriptionRepositoryBun) GetActiveAndUpcomingSubscriptions(ctx context.Context, companyID uuid.UUID) (*model.CompanySubscription, *model.CompanySubscription, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
		}
	}

	// Get Upcoming
	upcoming := &model.CompanySubscription{}
	if err := tx.NewSelect().
		Model(upcoming).
		Where("company_id = ?", companyID).
		Where("is_active = ?", true).
		Where("start_date > ?", time.Now().UTC()).
		Order("start_date ASC").
		Limit(1).
		Scan(ctx); err != nil {
		if err == sql.ErrNoRows {
			upcoming = nil
		} else {
			return nil, nil, err
		}
	}

	return active, upcoming, nil
}

func (r *CompanySubscriptionRepositoryBun) UpdateCompanyPlans(ctx context.Context) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	// 1. Update companies with active subscriptions
	// We select the latest/highest priority active subscription for each company
	// Ideally we should have a constraints to avoid overlap, but here we pick one.
	// We use a subquery or join.
	// Simple approach: Set Plan to Subscription's Plan where Subscription is Active and Valid.
	// "UPDATE companies SET current_plan = cs.plan_type FROM company_subscriptions cs WHERE cs.company_id = companies.id AND cs.is_active = true AND cs.start_date <= NOW() AND cs.end_date >= NOW()"
	if _, err := tx.NewRaw(`
		UPDATE companies 
		SET current_plan = cs.plan_type 
		FROM company_subscriptions cs 
		WHERE cs.company_id = companies.id 
		AND cs.is_active = true 
		AND cs.start_date <= NOW() 
		AND cs.end_date >= NOW()
	`).Exec(ctx); err != nil {
		return err
	}

	// 2. Mark expired subscriptions as inactive
	// This ensures is_active accurately reflects the subscription state
	if _, err := tx.NewRaw(`
		UPDATE company_subscriptions 
		SET is_active = false 
		WHERE is_active = true 
		AND end_date < NOW()
	`).Exec(ctx); err != nil {
		return err
	}

	// 3. Revert companies with NO active subscription to 'free'
	// "UPDATE companies SET current_plan = 'free' WHERE NOT EXISTS (SELECT 1 FROM company_subscriptions cs WHERE cs.company_id = companies.id AND cs.is_active = true AND cs.start_date <= NOW() AND cs.end_date >= NOW()) AND current_plan != 'free'"
	if _, err := tx.NewRaw(`
		UPDATE companies 
		SET current_plan = 'free' 
		WHERE NOT EXISTS (
			SELECT 1 FROM company_subscriptions cs 
			WHERE cs.company_id = companies.id 
			AND cs.is_active = true 
			AND cs.start_date <= NOW() 
			AND cs.end_date >= NOW()
		) 
		AND current_plan != 'free'
	`).Exec(ctx); err != nil {
		return err
	}

	// 4. Disable fiscal settings for companies on 'free' or 'basic' plans
	if _, err := tx.NewRaw(`
		UPDATE fiscal_settings 
		SET is_active = false 
		FROM companies c 
		WHERE fiscal_settings.company_id = c.id 
		AND c.current_plan IN ('free', 'basic') 
		AND fiscal_settings.is_active = true
	`).Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
