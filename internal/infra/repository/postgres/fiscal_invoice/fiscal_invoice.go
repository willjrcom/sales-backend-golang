package fiscalinvoicerepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type FiscalInvoiceRepository struct {
	db *bun.DB
}

func NewFiscalInvoiceRepository(db *bun.DB) *FiscalInvoiceRepository {
	return &FiscalInvoiceRepository{db: db}
}

func (r *FiscalInvoiceRepository) Create(ctx context.Context, invoice *model.FiscalInvoice) error {
	_, err := r.db.NewInsert().
		Model(invoice).
		Exec(ctx)
	return err
}

func (r *FiscalInvoiceRepository) Update(ctx context.Context, invoice *model.FiscalInvoice) error {
	_, err := r.db.NewUpdate().
		Model(invoice).
		WherePK().
		Exec(ctx)
	return err
}

func (r *FiscalInvoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FiscalInvoice, error) {
	invoice := &model.FiscalInvoice{}
	err := r.db.NewSelect().
		Model(invoice).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *FiscalInvoiceRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*model.FiscalInvoice, error) {
	invoice := &model.FiscalInvoice{}
	err := r.db.NewSelect().
		Model(invoice).
		Where("order_id = ?", orderID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *FiscalInvoiceRepository) GetByChaveAcesso(ctx context.Context, chaveAcesso string) (*model.FiscalInvoice, error) {
	invoice := &model.FiscalInvoice{}
	err := r.db.NewSelect().
		Model(invoice).
		Where("chave_acesso = ?", chaveAcesso).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *FiscalInvoiceRepository) List(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]*model.FiscalInvoice, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	var invoices []*model.FiscalInvoice
	total, err := r.db.NewSelect().
		Model(&invoices).
		Where("company_id = ?", companyID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (r *FiscalInvoiceRepository) GetNextNumero(ctx context.Context, companyID uuid.UUID, serie int) (int, error) {
	var maxNumero int
	err := r.db.NewSelect().
		Model((*model.FiscalInvoice)(nil)).
		ColumnExpr("COALESCE(MAX(numero), 0) as max_numero").
		Where("company_id = ?", companyID).
		Where("serie = ?", serie).
		Where("deleted_at IS NULL").
		Scan(ctx, &maxNumero)

	if err != nil {
		return 0, err
	}

	return maxNumero + 1, nil
}
