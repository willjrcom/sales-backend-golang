package addressentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewAddress(t *testing.T) {
   aca := &AddressCommonAttributes{
       Street: "Main", Number: "123", UF: "AA", Cep: "00000", AddressType: AddressTypeHouse, DeliveryTax: 2.5,
   }
   a := NewAddress(aca)
   assert.NotEqual(t, uuid.Nil, a.ID)
   assert.Equal(t, "Main", a.Street)
   assert.Equal(t, "123", a.Number)
   assert.Equal(t, "AA", a.UF)
   assert.Equal(t, "00000", a.Cep)
   assert.Equal(t, AddressTypeHouse, a.AddressType)
   assert.Equal(t, 2.5, a.DeliveryTax)
}