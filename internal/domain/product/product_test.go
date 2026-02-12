package productentity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

func TestNewProduct(t *testing.T) {
	cat := &ProductCategory{Entity: entity.NewEntity()}
	size := &Size{Entity: entity.NewEntity()}
	attrs := ProductCommonAttributes{SKU: "C", Name: "N", Category: cat, CategoryID: cat.ID, Size: size, SizeID: size.ID}
	p := NewProduct(attrs)
	assert.NotEqual(t, uuid.Nil, p.ID)
	assert.Equal(t, "C", p.SKU)
}

func TestFindSizeInCategory(t *testing.T) {
	cat := &ProductCategory{
		Entity: entity.NewEntity(),
		ProductCategoryCommonAttributes: ProductCategoryCommonAttributes{
			Sizes: []Size{{Entity: entity.NewEntity(), SizeCommonAttributes: SizeCommonAttributes{Name: "S", CategoryID: uuid.New()}}},
		},
	}
	ids := cat.Sizes[0].ID
	p := NewProduct(ProductCommonAttributes{Category: cat, CategoryID: cat.ID, SizeID: ids})
	found, err := p.FindSizeInCategory()
	assert.True(t, found)
	assert.NoError(t, err)
	p = NewProduct(ProductCommonAttributes{Category: cat, CategoryID: cat.ID, SizeID: uuid.New()})
	found, err = p.FindSizeInCategory()
	assert.False(t, found)
	assert.Error(t, err)
}
