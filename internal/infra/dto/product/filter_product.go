package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type FilterProductInput struct {
	Code       *string    `json:"code"`
	Name       *string    `json:"name"`
	CategoryID *uuid.UUID `json:"category_id"`
	SizeID     *uuid.UUID `json:"size_id"`
}

func (f *FilterProductInput) ToModel() *productentity.Product {
	product := &productentity.Product{}

	if f.Code != nil {
		product.Code = *f.Code
	}
	if f.Name != nil {
		product.Name = *f.Name
	}
	if f.CategoryID != nil {
		product.CategoryID = *f.CategoryID
	}
	if f.SizeID != nil {
		product.SizeID = *f.SizeID
	}

	return product
}
