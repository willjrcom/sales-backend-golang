package model

import (
	"context"

	"github.com/google/uuid"
)

type FiscalSettingsRepository interface {
	Create(ctx context.Context, fiscalSettings *FiscalSettings) error
	Update(ctx context.Context, fiscalSettings *FiscalSettings) error
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*FiscalSettings, error)
}
