package productcategoryproductdto

import (
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type UpdateProductInput struct {
	productentity.PatchProduct
}

func (p *UpdateProductInput) Validate(product *productentity.Product) error {
	if product.Code == "" {
		return ErrCodeRequired
	}
	if product.Name == "" {
		return ErrNameRequired
	}
	if product.Price < product.Cost {
		return ErrCostGreaterThanPrice
	}

	return nil
}

func (p *UpdateProductInput) UpdateModel(product *productentity.Product) (err error) {
	if p.Code != nil {
		product.Code = *p.Code
	}
	if p.Name != nil {
		product.Name = *p.Name
	}
	if p.Flavors != nil {
		product.Flavors = p.Flavors
	}
	if p.ImagePath != nil {
		product.ImagePath = p.ImagePath
	}
	if p.Description != nil {
		product.Description = *p.Description
	}
	if p.Price != nil {
		product.Price = *p.Price
	}
	if p.Cost != nil {
		product.Cost = *p.Cost
	}
	if p.IsAvailable != nil {
		product.IsAvailable = *p.IsAvailable
	}
	if p.CategoryID != nil {
		product.Category.ID = *p.CategoryID
	}
	if p.SizeID != nil {
		product.SizeID = *p.SizeID
	}

	if err = p.Validate(product); err != nil {
		return err
	}

	return nil
}
