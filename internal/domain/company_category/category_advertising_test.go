package companycategoryentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestCategoryToAdvertising(t *testing.T) {
   cid := uuid.New()
   aid := uuid.New()
   ca := CategoryToAdvertising{CompanyCategoryID: cid, AdvertisingID: aid}
   assert.Equal(t, cid, ca.CompanyCategoryID)
   assert.Equal(t, aid, ca.AdvertisingID)
}