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

	CompanyID         uuid.UUID `bun:"company_id,type:uuid,notnull"`
	Company           *Company  `bun:"rel:belongs-to,join:company_id=id"`
	PlanType          string    `bun:"plan_type,notnull"`
	StartDate         time.Time `bun:"start_date,notnull"`
	EndDate           time.Time `bun:"end_date,notnull"`
	IsActive          bool      `bun:"is_active,notnull"`
	IsCanceled        bool      `bun:"is_canceled,notnull,default:false"` // Renewal cancelled in MP
	PreapprovalID     *string   `bun:"preapproval_id"`                    // Unique index in DB
	ExternalReference *string   `bun:"external_reference"`
	Status            string    `bun:"status"`
}

func (c *CompanySubscription) ToDomain() *companyentity.CompanySubscription {
	return &companyentity.CompanySubscription{
		Entity:            c.Entity.ToDomain(),
		CompanyID:         c.CompanyID,
		PlanType:          companyentity.PlanType(c.PlanType),
		StartDate:         c.StartDate,
		EndDate:           c.EndDate,
		IsActive:          c.IsActive,
		IsCanceled:        c.IsCanceled,
		PreapprovalID:     c.PreapprovalID,
		ExternalReference: c.ExternalReference,
		Status:            c.Status,
	}
}

func (c *CompanySubscription) FromDomain(entity *companyentity.CompanySubscription) {
	c.Entity = entitymodel.FromDomain(entity.Entity)
	c.CompanyID = entity.CompanyID
	c.PlanType = string(entity.PlanType)
	c.StartDate = entity.StartDate
	c.EndDate = entity.EndDate
	c.IsActive = entity.IsActive
	c.IsCanceled = entity.IsCanceled
	c.PreapprovalID = entity.PreapprovalID
	c.ExternalReference = entity.ExternalReference
	c.Status = entity.Status
}
