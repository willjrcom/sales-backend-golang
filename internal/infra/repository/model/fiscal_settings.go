package model

import (
	"context"

	"github.com/google/uuid"
	fiscalsettingsentity "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_settings"
)

type FiscalSettingsRepository interface {
	Create(ctx context.Context, fiscalSettings *fiscalsettingsentity.FiscalSettings) error
	Update(ctx context.Context, fiscalSettings *fiscalsettingsentity.FiscalSettings) error
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*fiscalsettingsentity.FiscalSettings, error)
}
