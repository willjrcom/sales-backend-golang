package model

import (
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type ProductCategoryWithOrderProcess struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:product_categories"`
	ProductCategoryWithOrderProcessCommonAttributes
}

type ProductCategoryWithOrderProcessCommonAttributes struct {
	Name           string                        `bun:"name,notnull"`
	ImagePath      string                        `bun:"image_path"`
	UseProcessRule bool                          `bun:"use_process_rule"`
	IsAdditional   bool                          `bun:"is_additional"`
	IsComplement   bool                          `bun:"is_complement"`
	NeedPrint      bool                          `bun:"need_print"`
	PrinterName    string                        `bun:"printer_name"`
	ProcessRules   []ProcessRuleWithOrderProcess `bun:"rel:has-many,join:id=category_id"`
}

func (c *ProductCategoryWithOrderProcess) FromDomain(category *productentity.ProductCategory) {
	if category == nil {
		return
	}
	*c = ProductCategoryWithOrderProcess{
		Entity: entitymodel.FromDomain(category.Entity),
		ProductCategoryWithOrderProcessCommonAttributes: ProductCategoryWithOrderProcessCommonAttributes{
			Name:           category.Name,
			ImagePath:      category.ImagePath,
			UseProcessRule: category.UseProcessRule,
			IsAdditional:   category.IsAdditional,
			IsComplement:   category.IsComplement,
			NeedPrint:      category.NeedPrint,
			PrinterName:    category.PrinterName,
			ProcessRules:   []ProcessRuleWithOrderProcess{},
		},
	}

	for _, processRule := range category.ProcessRules {
		p := ProcessRuleWithOrderProcess{}
		p.FromDomain(&processRule)
		c.ProcessRules = append(c.ProcessRules, p)
	}
}

func (c *ProductCategoryWithOrderProcess) ToDomain() *productentity.ProductCategoryWithOrderProcess {
	if c == nil {
		return nil
	}
	category := &productentity.ProductCategoryWithOrderProcess{
		Entity: c.Entity.ToDomain(),
		ProductCategoryWithOrderProcessCommonAttributes: productentity.ProductCategoryWithOrderProcessCommonAttributes{
			Name:           c.Name,
			ImagePath:      c.ImagePath,
			UseProcessRule: c.UseProcessRule,
			IsAdditional:   c.IsAdditional,
			IsComplement:   c.IsComplement,
			NeedPrint:      c.NeedPrint,
			PrinterName:    c.PrinterName,
			ProcessRules:   []productentity.ProcessRuleWithOrderProcess{},
		},
	}

	for _, processRule := range c.ProcessRules {
		category.ProcessRules = append(category.ProcessRules, *processRule.ToDomain())
	}

	return category
}
