package model

import (
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type ProcessRuleWithOrderProcess struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:process_rules,alias:pr"`
	ProcessRuleCommonAttributes
	ProcessRuleWithOrderProcessCommonAttributes
}

type ProcessRuleWithOrderProcessCommonAttributes struct {
	TotalOrderProcessLate int `bun:"-"`
	TotalOrderQueue       int `bun:"-"`
}

func (p *ProcessRuleWithOrderProcess) FromDomain(processRule *productentity.ProcessRule) {
	if processRule == nil {
		return
	}
	*p = ProcessRuleWithOrderProcess{
		Entity: entitymodel.FromDomain(processRule.Entity),
		ProcessRuleCommonAttributes: ProcessRuleCommonAttributes{
			Name:              processRule.Name,
			Order:             processRule.Order,
			Description:       processRule.Description,
			ImagePath:         processRule.ImagePath,
			IdealTime:         processRule.IdealTime,
			ExperimentalError: processRule.ExperimentalError,
			CategoryID:        processRule.CategoryID,
		},
		ProcessRuleWithOrderProcessCommonAttributes: ProcessRuleWithOrderProcessCommonAttributes{
			TotalOrderProcessLate: 0,
			TotalOrderQueue:       0,
		},
	}
}

func (p *ProcessRuleWithOrderProcess) ToDomain() *productentity.ProcessRuleWithOrderProcess {
	if p == nil {
		return nil
	}
	return &productentity.ProcessRuleWithOrderProcess{
		Entity: p.Entity.ToDomain(),
		ProcessRuleCommonAttributes: productentity.ProcessRuleCommonAttributes{
			Name:              p.Name,
			Order:             p.Order,
			Description:       p.Description,
			ImagePath:         p.ImagePath,
			IdealTime:         p.IdealTime,
			ExperimentalError: p.ExperimentalError,
			CategoryID:        p.CategoryID,
		},
		ProcessRuleWithOrderProcessCommonAttributes: productentity.ProcessRuleWithOrderProcessCommonAttributes{
			TotalOrderProcessLate: p.TotalOrderProcessLate,
			TotalOrderQueue:       p.TotalOrderQueue,
		},
	}
}
