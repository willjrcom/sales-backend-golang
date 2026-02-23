package productcategorydto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"

	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
)

type CategoryDTO struct {
	ID                   uuid.UUID         `json:"id"`
	Name                 string            `json:"name"`
	ImagePath            string            `json:"image_path,omitempty"`
	NeedPrint            bool              `json:"need_print"`
	PrinterName          string            `json:"printer_name,omitempty"`
	UseProcessRule       bool              `json:"use_process_rule"`
	IsAdditional         bool              `json:"is_additional"`
	IsComplement         bool              `json:"is_complement"`
	AllowFractional      bool              `json:"allow_fractional"`
	IsActive             bool              `json:"is_active"`
	RemovableIngredients []string          `json:"removable_ingredients,omitempty"`
	Sizes                []sizedto.SizeDTO `json:"sizes,omitempty"`
	Products             []ProductDTO      `json:"products,omitempty"`
	AdditionalCategories []CategoryDTO     `json:"additional_categories,omitempty"`
	ComplementCategories []CategoryDTO     `json:"complement_categories,omitempty"`
	ProcessRules         []ProcessRuleDTO  `json:"process_rules,omitempty"`
}

func (c *CategoryDTO) FromDomain(category *productentity.ProductCategory) {
	if category == nil {
		return
	}
	*c = CategoryDTO{
		ID:                   category.ID,
		Name:                 category.Name,
		ImagePath:            category.ImagePath,
		NeedPrint:            category.NeedPrint,
		PrinterName:          category.PrinterName,
		UseProcessRule:       category.UseProcessRule,
		IsAdditional:         category.IsAdditional,
		IsComplement:         category.IsComplement,
		AllowFractional:      category.AllowFractional,
		IsActive:             category.IsActive,
		RemovableIngredients: category.RemovableIngredients,
		Sizes:                []sizedto.SizeDTO{},
		Products:             []ProductDTO{},
		AdditionalCategories: []CategoryDTO{},
		ComplementCategories: []CategoryDTO{},
		ProcessRules:         []ProcessRuleDTO{},
	}

	for _, processRule := range category.ProcessRules {
		p := ProcessRuleDTO{}
		p.FromDomain(&processRule)
		c.ProcessRules = append(c.ProcessRules, p)
	}

	for _, size := range category.Sizes {
		s := sizedto.SizeDTO{}
		s.FromDomain(&size)
		c.Sizes = append(c.Sizes, s)
	}

	for _, product := range category.Products {
		p := ProductDTO{}
		p.FromDomain(&product)
		c.Products = append(c.Products, p)
	}

	for _, category := range category.AdditionalCategories {
		a := CategoryDTO{}
		a.FromDomain(&category)
		c.AdditionalCategories = append(c.AdditionalCategories, a)
	}

	for _, category := range category.ComplementCategories {
		a := CategoryDTO{}
		a.FromDomain(&category)
		c.ComplementCategories = append(c.ComplementCategories, a)
	}

	if len(category.ProcessRules) == 0 {
		c.ProcessRules = nil
	}
	if len(category.Sizes) == 0 {
		c.Sizes = nil
	}
	if len(category.Products) == 0 {
		c.Products = nil
	}
	if len(category.AdditionalCategories) == 0 {
		c.AdditionalCategories = nil
	}
	if len(category.ComplementCategories) == 0 {
		c.ComplementCategories = nil
	}
}
