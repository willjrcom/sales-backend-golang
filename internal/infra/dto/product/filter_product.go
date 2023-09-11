package productdto

import (
	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type FilterProductInput struct {
	Code       *string
	Name       *string
	CategoryID *uuid.UUID
	Size       *string
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
		product.Category.ID = *f.CategoryID
	}
	if f.Size != nil {
		product.Size = *f.Size
	}

	return product
}
