package productcategorydto

import (
	"errors"

	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameIsEmpty = errors.New("name is empty")
)

type CategoryCreateDTO struct {
	Name                 string   `json:"name"`
	ImagePath            string   `json:"image_path"`
	NeedPrint            bool     `json:"need_print"`
	PrinterName          string   `json:"printer_name"`
	UseProcessRule       bool     `json:"use_process_rule"`
	RemovableIngredients []string `json:"removable_ingredients"`
	IsAdditional         bool     `json:"is_additional"`
	IsComplement         bool     `json:"is_complement"`
	AllowFractional      bool     `json:"allow_fractional"`
	IsActive             *bool    `json:"is_active"`
}

func (c *CategoryCreateDTO) validate() error {
	if c.Name == "" {
		return ErrNameIsEmpty
	}

	return nil
}

func (c *CategoryCreateDTO) ToDomain() (*productentity.ProductCategory, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	isActive := true
	if c.IsActive != nil {
		isActive = *c.IsActive
	}

	categoryCommonAttributes := productentity.ProductCategoryCommonAttributes{
		Name:                 c.Name,
		ImagePath:            c.ImagePath,
		NeedPrint:            c.NeedPrint,
		PrinterName:          c.PrinterName,
		UseProcessRule:       c.UseProcessRule,
		RemovableIngredients: c.RemovableIngredients,
		IsAdditional:         c.IsAdditional,
		IsComplement:         c.IsComplement,
		IsActive:             isActive,
		AllowFractional:      c.AllowFractional,
	}

	return productentity.NewProductCategory(categoryCommonAttributes), nil
}
