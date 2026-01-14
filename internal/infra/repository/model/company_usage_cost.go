package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanyUsageCost struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:company_usage_costs"`

	CompanyID    uuid.UUID       `bun:"company_id,type:uuid,notnull"`
	CostType     string          `bun:"cost_type,notnull"`
	Description  string          `bun:"description"`
	Amount       decimal.Decimal `bun:"amount,type:decimal(19,4),notnull"`
	ReferenceID  *uuid.UUID      `bun:"reference_id,type:uuid"`
	BillingMonth int             `bun:"billing_month,notnull"`
	BillingYear  int             `bun:"billing_year,notnull"`
}

func (c *CompanyUsageCost) FromDomain(cost *companyentity.CompanyUsageCost) {
	if cost == nil {
		return
	}
	*c = CompanyUsageCost{
		Entity:       entitymodel.FromDomain(cost.Entity),
		CompanyID:    cost.CompanyID,
		CostType:     string(cost.CostType),
		Description:  cost.Description,
		Amount:       cost.Amount,
		ReferenceID:  cost.ReferenceID,
		BillingMonth: cost.BillingMonth,
		BillingYear:  cost.BillingYear,
	}
}

func (c *CompanyUsageCost) ToDomain() *companyentity.CompanyUsageCost {
	if c == nil {
		return nil
	}
	return &companyentity.CompanyUsageCost{
		Entity:       c.Entity.ToDomain(),
		CompanyID:    c.CompanyID,
		CostType:     companyentity.CostType(c.CostType),
		Description:  c.Description,
		Amount:       c.Amount,
		ReferenceID:  c.ReferenceID,
		BillingMonth: c.BillingMonth,
		BillingYear:  c.BillingYear,
	}
}

// CompanyUsageCostRepository defines the interface for company usage cost operations
type CompanyUsageCostRepository interface {
	Create(ctx context.Context, cost *CompanyUsageCost) error
	GetByID(ctx context.Context, id uuid.UUID) (*CompanyUsageCost, error)
	GetMonthlyCosts(ctx context.Context, companyID uuid.UUID, month, year int) ([]*CompanyUsageCost, error)
	GetCostsByType(ctx context.Context, companyID uuid.UUID, costType string, month, year int) ([]*CompanyUsageCost, error)
	GetTotalByMonth(ctx context.Context, companyID uuid.UUID, month, year int) (decimal.Decimal, error)
	GetByReferenceID(ctx context.Context, referenceID uuid.UUID) (*CompanyUsageCost, error)
}
