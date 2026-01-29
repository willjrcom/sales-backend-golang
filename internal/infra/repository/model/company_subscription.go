package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanySubscription struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:company_subscriptions"`

	CompanyID uuid.UUID              `bun:"company_id,type:uuid,notnull"`
	Company   *Company               `bun:"rel:belongs-to,join:company_id=id"`
	PaymentID *uuid.UUID             `bun:"payment_id,type:uuid"` // Nullable
	Payment   *CompanyPayment        `bun:"rel:belongs-to,join:payment_id=id"`
	PlanType  companyentity.PlanType `bun:"plan_type,notnull"`
	StartDate time.Time              `bun:"start_date,notnull"`
	EndDate   time.Time              `bun:"end_date,notnull"`
	IsActive  bool                   `bun:"is_active,notnull"`
}

func (c *CompanySubscription) ToDomain() *companyentity.CompanySubscription {
	return &companyentity.CompanySubscription{
		Entity:    c.Entity.ToDomain(),
		CompanyID: c.CompanyID,
		PaymentID: c.PaymentID,
		PlanType:  c.PlanType,
		StartDate: c.StartDate,
		EndDate:   c.EndDate,
		IsActive:  c.IsActive,
	}
}

func (c *CompanySubscription) FromDomain(entity *companyentity.CompanySubscription) {
	c.Entity = entitymodel.FromDomain(entity.Entity)
	c.CompanyID = entity.CompanyID
	c.PaymentID = entity.PaymentID
	c.PlanType = entity.PlanType
	c.StartDate = entity.StartDate
	c.EndDate = entity.EndDate
	c.IsActive = entity.IsActive
}
