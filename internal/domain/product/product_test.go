package productentity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

func TestNewProduct(t *testing.T) {
	cat := &ProductCategory{Entity: entity.NewEntity()}
	attrs := ProductCommonAttributes{SKU: "C", Name: "N", Category: cat, CategoryID: cat.ID}
	p := NewProduct(attrs)
	assert.NotEqual(t, uuid.Nil, p.ID)
	assert.Equal(t, "C", p.SKU)
}

func TestAddVariation(t *testing.T) {
	cat := &ProductCategory{Entity: entity.NewEntity()}
	attrs := ProductCommonAttributes{SKU: "C", Name: "N", Category: cat, CategoryID: cat.ID}
	p := NewProduct(attrs)

	v := ProductVariation{
		Entity:    entity.NewEntity(),
		ProductID: p.ID,
		SizeID:    uuid.New(),
	}

	p.AddVariation(v)
	assert.Equal(t, 1, len(p.Variations))
	assert.Equal(t, v.ID, p.Variations[0].ID)
}
