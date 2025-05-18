package productentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
   "github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

func TestNewQuantity(t *testing.T) {
   q := NewQuantity(QuantityCommonAttributes{Quantity: 2, CategoryID: uuid.New()})
   assert.NotEqual(t, uuid.Nil, q.ID)
   assert.Equal(t, 2.0, q.Quantity)
}

func TestValidateDuplicateQuantities(t *testing.T) {
   id := uuid.New()
   q1 := *NewQuantity(QuantityCommonAttributes{Quantity: 1, CategoryID: id})
   q2 := *NewQuantity(QuantityCommonAttributes{Quantity: 2, CategoryID: id})
   assert.NoError(t, ValidateDuplicateQuantities(3, []Quantity{q1, q2}))
   assert.Equal(t, ErrQuantityAlreadyExists, ValidateDuplicateQuantities(2, []Quantity{q1, q2}))
}

func TestValidateUpdateQuantity(t *testing.T) {
   id1 := uuid.New()
   id2 := uuid.New()
   qs := []Quantity{
       {Entity: entity.NewEntity(), QuantityCommonAttributes: QuantityCommonAttributes{Quantity: 1, CategoryID: uuid.New(),},},
       {Entity: entity.NewEntity(), QuantityCommonAttributes: QuantityCommonAttributes{Quantity: 2, CategoryID: uuid.New(),},},
   }
   qs[0].ID = id1
   qs[1].ID = id2
   newQ := &Quantity{Entity: entity.NewEntity(), QuantityCommonAttributes: QuantityCommonAttributes{Quantity: 1, CategoryID: uuid.New()}}
   newQ.ID = id2
   assert.Equal(t, ErrQuantityAlreadyExists, ValidateUpdateQuantity(newQ, qs))
   newQ.Quantity = 3
   assert.NoError(t, ValidateUpdateQuantity(newQ, qs))
}