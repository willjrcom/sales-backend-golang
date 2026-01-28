package companyrepositorybun

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyPaymentRepositoryBun struct {
	db *bun.DB
}

func NewCompanyPaymentRepositoryBun(db *bun.DB) *CompanyPaymentRepositoryBun {
	return &CompanyPaymentRepositoryBun{db: db}
}

func (r *CompanyPaymentRepositoryBun) CreateCompanyPayment(ctx context.Context, payment *model.CompanyPayment) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(payment).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyPaymentRepositoryBun) UpdateCompanyPayment(ctx context.Context, payment *model.CompanyPayment) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(payment).WherePK().Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyPaymentRepositoryBun) GetCompanyPaymentByID(ctx context.Context, id uuid.UUID) (*model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	payment := &model.CompanyPayment{}
	if err := tx.NewSelect().
		Model(payment).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *CompanyPaymentRepositoryBun) ListCompanyPayments(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]model.CompanyPayment, int, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	total, err := tx.NewSelect().
		Model((*model.CompanyPayment)(nil)).
		Where("company_id = ?", companyID).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		if err := tx.Commit(); err != nil {
			return nil, 0, err
		}
		return []model.CompanyPayment{}, 0, nil
	}

	payments := make([]model.CompanyPayment, 0, perPage)

	if err := tx.NewSelect().
		Model(&payments).
		Where("company_id = ?", companyID).
		Order("paid_at DESC").
		Order("created_at DESC").
		Limit(perPage).
		Offset(page * perPage).
		Scan(ctx); err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *CompanyPaymentRepositoryBun) ListOverduePayments(ctx context.Context, cutoffDate time.Time) ([]model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	var payments []model.CompanyPayment
	if err := tx.NewSelect().
		Model(&payments).
		Where("status = ?", "pending").
		Where("is_mandatory = ?", true).
		Where("expires_at < ?", cutoffDate).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *CompanyPaymentRepositoryBun) ListPendingMandatoryPayments(ctx context.Context, companyID uuid.UUID) ([]model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	var payments []model.CompanyPayment
	if err := tx.NewSelect().
		Model(&payments).
		Where("company_id = ?", companyID).
		Where("status = ?", "pending").
		Where("is_mandatory = ?", true).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payments, nil
}
