package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type ProcessRule struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:process_rules,alias:pr"`
	ProcessRuleCommonAttributes
}

type ProcessRuleCommonAttributes struct {
	Name        string           `bun:"name,notnull"`
	Order       int8             `bun:"order,notnull"`
	Description string           `bun:"description"`
	ImagePath   *string          `bun:"image_path"`
	IdealTime   time.Duration    `bun:"ideal_time,notnull"`
	CategoryID  uuid.UUID        `bun:"column:category_id,type:uuid,notnull"`
	Category    *ProductCategory `bun:"rel:belongs-to"`
	IsActive    bool             `bun:"is_active,notnull"`
}

func (p *ProcessRule) FromDomain(processRule *productentity.ProcessRule) {
	if processRule == nil {
		return
	}
	*p = ProcessRule{
		Entity: entitymodel.FromDomain(processRule.Entity),
		ProcessRuleCommonAttributes: ProcessRuleCommonAttributes{
			Name:        processRule.Name,
			Order:       processRule.Order,
			Description: processRule.Description,
			ImagePath:   processRule.ImagePath,
			IdealTime:   processRule.IdealTime,
			CategoryID:  processRule.CategoryID,
			Category:    &ProductCategory{},
			IsActive:    processRule.IsActive,
		},
	}

	if processRule.Category != nil {
		p.Category.FromDomain(processRule.Category)
	}
}

func (p *ProcessRule) ToDomain() *productentity.ProcessRule {
	if p == nil {
		return nil
	}
	processRule := &productentity.ProcessRule{
		Entity: p.Entity.ToDomain(),
		ProcessRuleCommonAttributes: productentity.ProcessRuleCommonAttributes{
			Name:        p.Name,
			Order:       p.Order,
			Description: p.Description,
			ImagePath:   p.ImagePath,
			IdealTime:   p.IdealTime,
			CategoryID:  p.CategoryID,
			Category:    p.Category.ToDomain(),
			IsActive:    p.IsActive,
		},
	}

	return processRule
}
