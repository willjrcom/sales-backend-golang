package sponsorentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewSponsor(t *testing.T) {
   attr := SponsorCommonAttributes{Name: "N", CNPJ: "C", Email: "e", Contacts: []string{"c"}, Address: nil}
   s := NewSponsor(attr)
   assert.NotEqual(t, uuid.Nil, s.ID)
   assert.Equal(t, "N", s.Name)
   assert.Equal(t, "C", s.CNPJ)
   assert.Equal(t, "e", s.Email)
   assert.Equal(t, []string{"c"}, s.Contacts)
}