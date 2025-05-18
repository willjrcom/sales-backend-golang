package tableentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewPlace(t *testing.T) {
   pca := PlaceCommonAttributes{Name: "Place", ImagePath: nil, IsAvailable: true}
   p := NewPlace(pca)
   assert.NotEqual(t, uuid.Nil, p.ID)
   assert.Equal(t, "Place", p.Name)
}

func TestNewPlaceToTable(t *testing.T) {
   pid := uuid.New()
   tid := uuid.New()
   pt := NewPlaceToTable(pid, tid, 1, 2)
   assert.Equal(t, pid, pt.PlaceID)
   assert.Equal(t, tid, pt.TableID)
   assert.Equal(t, 1, pt.Column)
   assert.Equal(t, 2, pt.Row)
}