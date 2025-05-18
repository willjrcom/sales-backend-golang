package companycategoryentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewCategory(t *testing.T) {
   cca := CompanyCategoryCommonAttributes{Name: "Test", ImagePath: "img"}
   c := NewCategory(cca)
   assert.NotEqual(t, uuid.Nil, c.ID)
   assert.Equal(t, "Test", c.Name)
   assert.Equal(t, "img", c.ImagePath)
}