package productentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
   "github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

func TestNewSize(t *testing.T) {
   s := NewSize(SizeCommonAttributes{Name: "S", CategoryID: uuid.New()})
   assert.NotEqual(t, uuid.Nil, s.ID)
   assert.Equal(t, "S", s.Name)
}

func TestValidateDuplicateSizes(t *testing.T) {
   id := uuid.New()
   sz1 := *NewSize(SizeCommonAttributes{Name: "X", CategoryID: id})
   sz2 := *NewSize(SizeCommonAttributes{Name: "Y", CategoryID: id})
   assert.NoError(t, ValidateDuplicateSizes("Z", []Size{sz1, sz2}))
   assert.Equal(t, ErrSizeAlreadyExists, ValidateDuplicateSizes("X", []Size{sz1, sz2}))
}

func TestValidateUpdateSize(t *testing.T) {
   id1 := uuid.New()
   id2 := uuid.New()
   ss := []Size{
       {Entity: entity.NewEntity(), SizeCommonAttributes: SizeCommonAttributes{Name: "A", CategoryID: uuid.New()}},
       {Entity: entity.NewEntity(), SizeCommonAttributes: SizeCommonAttributes{Name: "B", CategoryID: uuid.New()}},
   }
   ss[0].ID = id1
   ss[1].ID = id2
   newS := &Size{Entity: entity.NewEntity(), SizeCommonAttributes: SizeCommonAttributes{Name: "A", CategoryID: uuid.New()}}
   newS.ID = id2
   assert.Equal(t, ErrSizeAlreadyExists, ValidateUpdateSize(newS, ss))
   newS.Name = "C"
   assert.NoError(t, ValidateUpdateSize(newS, ss))
}