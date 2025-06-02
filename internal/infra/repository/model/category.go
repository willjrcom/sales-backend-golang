package model

import (
	"github.com/uptrace/bun"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type ProductCategory struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:product_categories,alias:category"`
	ProductCategoryCommonAttributes
}

type ProductCategoryCommonAttributes struct {
	Name                 string            `bun:"name,notnull"`
	ImagePath            string            `bun:"image_path"`
	NeedPrint            bool              `bun:"need_print,notnull"`
	UseProcessRule       bool              `bun:"use_process_rule,notnull"`
	RemovableIngredients []string          `bun:"removable_ingredients,type:jsonb"`
	Sizes                []Size            `bun:"rel:has-many,join:id=category_id"`
	Quantities           []Quantity        `bun:"rel:has-many,join:id=category_id"`
	Products             []Product         `bun:"rel:has-many,join:id=category_id"`
	ProcessRules         []ProcessRule     `bun:"rel:has-many,join:id=category_id"`
	IsAdditional         bool              `bun:"is_additional"`
	IsComplement         bool              `bun:"is_complement"`
	AdditionalCategories []ProductCategory `bun:"m2m:product_category_to_additional,join:Category=AdditionalCategory"`
	ComplementCategories []ProductCategory `bun:"m2m:product_category_to_complement,join:Category=ComplementCategory"`
}

func (c *ProductCategory) FromDomain(category *productentity.ProductCategory) {
	if category == nil {
		return
	}
	*c = ProductCategory{
		Entity: entitymodel.FromDomain(category.Entity),
		ProductCategoryCommonAttributes: ProductCategoryCommonAttributes{
			Name:                 category.Name,
			ImagePath:            category.ImagePath,
			NeedPrint:            category.NeedPrint,
			UseProcessRule:       category.UseProcessRule,
			RemovableIngredients: category.RemovableIngredients,
			Sizes:                []Size{},
			Quantities:           []Quantity{},
			Products:             []Product{},
			ProcessRules:         []ProcessRule{},
			IsAdditional:         category.IsAdditional,
			IsComplement:         category.IsComplement,
			AdditionalCategories: []ProductCategory{},
			ComplementCategories: []ProductCategory{},
		},
	}

	for _, size := range category.Sizes {
		s := Size{}
		s.FromDomain(&size)
		c.Sizes = append(c.Sizes, s)
	}

	for _, quantity := range category.Quantities {
		q := Quantity{}
		q.FromDomain(&quantity)
		c.Quantities = append(c.Quantities, q)
	}

	for _, product := range category.Products {
		p := Product{}
		p.FromDomain(&product)
		c.Products = append(c.Products, p)
	}

	for _, processRule := range category.ProcessRules {
		p := ProcessRule{}
		p.FromDomain(&processRule)
		c.ProcessRules = append(c.ProcessRules, p)
	}

	for _, additionalCategory := range category.AdditionalCategories {
		a := ProductCategory{}
		a.FromDomain(&additionalCategory)
		c.AdditionalCategories = append(c.AdditionalCategories, a)
	}

	for _, complementCategory := range category.ComplementCategories {
		a := ProductCategory{}
		a.FromDomain(&complementCategory)
		c.ComplementCategories = append(c.ComplementCategories, a)
	}
}

func (c *ProductCategory) ToDomain() *productentity.ProductCategory {
	if c == nil {
		return nil
	}
	category := &productentity.ProductCategory{
		Entity: c.Entity.ToDomain(),
		ProductCategoryCommonAttributes: productentity.ProductCategoryCommonAttributes{
			Name:                 c.Name,
			ImagePath:            c.ImagePath,
			NeedPrint:            c.NeedPrint,
			UseProcessRule:       c.UseProcessRule,
			RemovableIngredients: c.RemovableIngredients,
			Sizes:                []productentity.Size{},
			Quantities:           []productentity.Quantity{},
			Products:             []productentity.Product{},
			ProcessRules:         []productentity.ProcessRule{},
			IsAdditional:         c.IsAdditional,
			IsComplement:         c.IsComplement,
			AdditionalCategories: []productentity.ProductCategory{},
			ComplementCategories: []productentity.ProductCategory{},
		},
	}

	for _, size := range c.Sizes {
		category.Sizes = append(category.Sizes, *size.ToDomain())
	}

	for _, quantity := range c.Quantities {
		category.Quantities = append(category.Quantities, *quantity.ToDomain())
	}

	for _, product := range c.Products {
		category.Products = append(category.Products, *product.ToDomain())
	}

	for _, processRule := range c.ProcessRules {
		category.ProcessRules = append(category.ProcessRules, *processRule.ToDomain())
	}

	for _, additionalCategory := range c.AdditionalCategories {
		category.AdditionalCategories = append(category.AdditionalCategories, *additionalCategory.ToDomain())
	}

	for _, complementCategory := range c.ComplementCategories {
		category.ComplementCategories = append(category.ComplementCategories, *complementCategory.ToDomain())
	}

	return category
}
