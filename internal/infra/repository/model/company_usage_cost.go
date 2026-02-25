package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanyUsageCost struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:company_usage_costs"`

	CompanyID   uuid.UUID        `bun:"company_id,type:uuid,notnull"`
	CostType    string           `bun:"cost_type,notnull"`
	Description string           `bun:"description"`
	Amount      *decimal.Decimal `bun:"amount,type:decimal(19,4),notnull"`
	Status      string           `bun:"status,notnull"`
	ReferenceID *uuid.UUID       `bun:"reference_id,type:uuid"`
	PaymentID   *uuid.UUID       `bun:"payment_id,type:uuid"`
}

func (c *CompanyUsageCost) FromDomain(cost *companyentity.CompanyUsageCost) {
	if cost == nil {
		return
	}
	*c = CompanyUsageCost{
		Entity:      entitymodel.FromDomain(cost.Entity),
		CompanyID:   cost.CompanyID,
		CostType:    string(cost.CostType),
		Description: cost.Description,
		Amount:      &cost.Amount,
		Status:      string(cost.Status),
		ReferenceID: cost.ReferenceID,
		PaymentID:   cost.PaymentID,
	}
}

func (c *CompanyUsageCost) ToDomain() *companyentity.CompanyUsageCost {
	if c == nil {
		return nil
	}
	return &companyentity.CompanyUsageCost{
		Entity:      c.Entity.ToDomain(),
		CompanyID:   c.CompanyID,
		CostType:    companyentity.CostType(c.CostType),
		Description: c.Description,
		Amount:      c.GetAmount(),
		Status:      companyentity.CostStatus(c.Status),
		ReferenceID: c.ReferenceID,
		PaymentID:   c.PaymentID,
	}
}

func (c *CompanyUsageCost) GetAmount() decimal.Decimal {
	if c.Amount == nil {
		return decimal.Zero
	}
	return *c.Amount
}
