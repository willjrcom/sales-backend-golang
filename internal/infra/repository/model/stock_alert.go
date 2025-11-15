package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	stockentity "github.com/willjrcom/sales-backend-go/internal/domain/stock"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

// StockAlert model
type StockAlert struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:stock_alerts,alias:stock_alert"`
	StockAlertCommonAttributes
}

type StockAlertCommonAttributes struct {
	StockID    uuid.UUID  `bun:"stock_id,type:uuid,notnull"`
	Type       string     `bun:"type,notnull"`
	Message    string     `bun:"message,notnull"`
	IsResolved bool       `bun:"is_resolved,notnull"`
	ResolvedAt *time.Time `bun:"resolved_at"`
	ResolvedBy *uuid.UUID `bun:"resolved_by,type:uuid"`
}

// FromDomain converte domain para model
func (sa *StockAlert) FromDomain(alert *stockentity.StockAlert) {
	if alert == nil {
		return
	}
	*sa = StockAlert{
		Entity: entitymodel.FromDomain(alert.Entity),
		StockAlertCommonAttributes: StockAlertCommonAttributes{
			StockID:    alert.StockID,
			Type:       string(alert.Type),
			Message:    alert.Message,
			IsResolved: alert.IsResolved,
			ResolvedAt: alert.ResolvedAt,
			ResolvedBy: alert.ResolvedBy,
		},
	}
}

// ToDomain converte model para domain
func (sa *StockAlert) ToDomain() *stockentity.StockAlert {
	if sa == nil {
		return nil
	}
	return &stockentity.StockAlert{
		Entity: sa.Entity.ToDomain(),
		StockAlertCommonAttributes: stockentity.StockAlertCommonAttributes{
			StockID:    sa.StockID,
			Type:       stockentity.AlertType(sa.Type),
			Message:    sa.Message,
			IsResolved: sa.IsResolved,
			ResolvedAt: sa.ResolvedAt,
			ResolvedBy: sa.ResolvedBy,
		},
	}
}
