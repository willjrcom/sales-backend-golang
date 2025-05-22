package productcategorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	processruledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_rule"
)

type CategoryWithOrderProcessDTO struct {
	ID             uuid.UUID                                       `json:"id"`
	Name           string                                          `json:"name"`
	ImagePath      string                                          `json:"image_path,omitempty"`
	NeedPrint      bool                                            `json:"need_print"`
	UseProcessRule bool                                            `json:"use_process_rule"`
	IsAdditional   bool                                            `json:"is_additional"`
	IsComplement   bool                                            `json:"is_complement"`
	ProcessRules   []processruledto.ProcessRuleWithOrderProcessDTO `json:"process_rules,omitempty"`
}

func (c *CategoryWithOrderProcessDTO) FromDomain(category *productentity.ProductCategoryWithOrderProcess) {
	if category == nil {
		return
	}
	*c = CategoryWithOrderProcessDTO{
		ID:             category.ID,
		Name:           category.Name,
		ImagePath:      category.ImagePath,
		UseProcessRule: category.UseProcessRule,
		IsAdditional:   category.IsAdditional,
		IsComplement:   category.IsComplement,
		NeedPrint:      category.NeedPrint,
		ProcessRules:   []processruledto.ProcessRuleWithOrderProcessDTO{},
	}

	for _, processRule := range category.ProcessRules {
		p := processruledto.ProcessRuleWithOrderProcessDTO{}
		p.FromDomain(&processRule)
		c.ProcessRules = append(c.ProcessRules, p)
	}

	if len(category.ProcessRules) == 0 {
		c.ProcessRules = nil
	}
}
