package fiscalinvoicerepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type FiscalInvoiceRepository struct {
	db *bun.DB
}

func NewFiscalInvoiceRepository(db *bun.DB) *FiscalInvoiceRepository {
	return &FiscalInvoiceRepository{db: db}
}

func (r *FiscalInvoiceRepository) Create(ctx context.Context, invoice *model.FiscalInvoice) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().
		Model(invoice).
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *FiscalInvoiceRepository) Update(ctx context.Context, invoice *model.FiscalInvoice) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().
		Model(invoice).
		WherePK().
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *FiscalInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FiscalInvoice, error) {
	invoice := &model.FiscalInvoice{}
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(invoice).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *FiscalInvoiceRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*model.FiscalInvoice, error) {
	invoice := &model.FiscalInvoice{}
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(invoice).
		Where("order_id = ?", orderID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *FiscalInvoiceRepository) GetByAccessKey(ctx context.Context, accessKey string) (*model.FiscalInvoice, error) {
	invoice := &model.FiscalInvoice{}
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model(invoice).
		Where("access_key = ?", accessKey).
		Where("deleted_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *FiscalInvoiceRepository) List(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]*model.FiscalInvoice, int, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	invoices := []*model.FiscalInvoice{}
	total, err := tx.NewSelect().
		Model(&invoices).
		Where("company_id = ?", companyID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(perPage).
		Offset(page * perPage).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

func (r *FiscalInvoiceRepository) GetNextNumber(ctx context.Context, companyID uuid.UUID, series int) (int, error) {
	var maxNumber int
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return 0, err
	}
	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().
		Model((*model.FiscalInvoice)(nil)).
		ColumnExpr("COALESCE(MAX(number), 0) as max_number").
		Where("company_id = ?", companyID).
		Where("series = ?", series).
		Where("deleted_at IS NULL").
		Scan(ctx, &maxNumber); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return maxNumber + 1, nil
}
