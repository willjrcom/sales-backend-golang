package advertisingentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewAdvertising(t *testing.T) {
   aca := AdvertisingCommonAttributes{Name: "Name", ImagePath: "path"}
   a := NewAdvertising(aca)
   assert.NotEqual(t, uuid.Nil, a.ID)
   assert.Equal(t, "Name", a.Name)
   assert.Equal(t, "path", a.ImagePath)
}