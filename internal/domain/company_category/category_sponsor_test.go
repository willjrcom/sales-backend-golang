package companycategoryentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestCategoryToSponsor(t *testing.T) {
   cid := uuid.New()
   sid := uuid.New()
   cs := CategoryToSponsor{CompanyCategoryID: cid, SponsorID: sid}
   assert.Equal(t, cid, cs.CompanyCategoryID)
   assert.Equal(t, sid, cs.SponsorID)
}