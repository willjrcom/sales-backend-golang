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

func (r *CompanyPaymentRepositoryBun) ListCompanyPayments(ctx context.Context, companyID uuid.UUID, page, perPage, month, year int) ([]model.CompanyPayment, int, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	q := tx.NewSelect().
		Model((*model.CompanyPayment)(nil)).
		Where("company_id = ?", companyID)

	if month > 0 && year > 0 {
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)
		q.Where("created_at >= ?", startDate).Where("created_at < ?", endDate)
	}

	total, err := q.Count(ctx)
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

	qList := tx.NewSelect().
		Model(&payments).
		Where("company_id = ?", companyID)

	if month > 0 && year > 0 {
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)
		qList.Where("created_at >= ?", startDate).Where("created_at < ?", endDate)
	}

	if err := qList.
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

func (r *CompanyPaymentRepositoryBun) GetCompanyPaymentByProviderID(ctx context.Context, providerPaymentID string) (*model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	payment := &model.CompanyPayment{}
	query := tx.NewSelect().
		Model(payment).
		Where("provider_payment_id = ?", providerPaymentID).
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *CompanyPaymentRepositoryBun) GetPendingPaymentByExternalReference(ctx context.Context, externalReference string) (*model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	payment := &model.CompanyPayment{}
	query := tx.NewSelect().
		Model(payment).
		Where("external_reference = ?", externalReference).
		Where("status = ?", "pending").
		Order("created_at DESC").
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *CompanyPaymentRepositoryBun) GetCompanyPaymentByExternalReference(ctx context.Context, externalReference string) (*model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	payment := &model.CompanyPayment{}
	query := tx.NewSelect().
		Model(payment).
		Where("external_reference = ?", externalReference).
		Order("created_at DESC").
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *CompanyPaymentRepositoryBun) GetLastApprovedPaymentByExternalReferencePrefix(ctx context.Context, externalReferencePrefix string) (*model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	payment := &model.CompanyPayment{}
	query := tx.NewSelect().
		Model(payment).
		Where("external_reference LIKE ?", externalReferencePrefix+"%").
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("status = ?", "approved").
				WhereOr("status = ?", "paid")
		}).
		Order("created_at DESC").
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *CompanyPaymentRepositoryBun) GetLastPaymentByExternalReferencePrefix(ctx context.Context, externalReferencePrefix string) (*model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	payment := &model.CompanyPayment{}
	query := tx.NewSelect().
		Model(payment).
		Where("external_reference LIKE ?", externalReferencePrefix+"%").
		Order("created_at DESC").
		Limit(1)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *CompanyPaymentRepositoryBun) ListOverduePaymentsByCompany(ctx context.Context, companyID uuid.UUID, cutoffDate time.Time) ([]model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	var payments []model.CompanyPayment
	query := tx.NewSelect().
		Model(&payments).
		Where("company_id = ?", companyID).
		Where("status = ?", "pending").
		Where("is_mandatory = ?", true).
		Where("expires_at < ?", cutoffDate)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *CompanyPaymentRepositoryBun) ListExpiredOptionalPayments(ctx context.Context) ([]model.CompanyPayment, error) {
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
		Where("is_mandatory = ?", false).
		Where("expires_at < ?", time.Now().UTC()).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return payments, nil
}
