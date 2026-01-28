package fiscalsettingsrepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type FiscalSettingsRepositoryBun struct {
	db *bun.DB
}

func NewFiscalSettingsRepositoryBun(db *bun.DB) model.FiscalSettingsRepository {
	return &FiscalSettingsRepositoryBun{db: db}
}

func (r *FiscalSettingsRepositoryBun) Create(ctx context.Context, settings *model.FiscalSettings) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(settings).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *FiscalSettingsRepositoryBun) Update(ctx context.Context, settings *model.FiscalSettings) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(settings).WherePK().Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *FiscalSettingsRepositoryBun) GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*model.FiscalSettings, error) {
	m := &model.FiscalSettings{}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	err = tx.NewSelect().Model(m).Where("company_id = ?", companyID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return m, nil
}
