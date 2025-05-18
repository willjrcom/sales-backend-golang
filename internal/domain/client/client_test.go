package cliententity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
   personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
   addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

func TestNewClient(t *testing.T) {
   pc := &personentity.PersonCommonAttributes{Name: "N", Email: "e", Cpf: "c"}
   p := personentity.NewPerson(pc)
   c := NewClient(p)
   assert.NotEqual(t, uuid.Nil, c.ID)
   assert.Equal(t, p.PersonCommonAttributes, c.Person.PersonCommonAttributes)
}

func TestAddContactAndAddress(t *testing.T) {
   pc := &personentity.PersonCommonAttributes{Name: "N"}
   p := personentity.NewPerson(pc)
   c := NewClient(p)
   contactAttr := &personentity.ContactCommonAttributes{Ddd: "012", Number: "1111", Type: personentity.ContactTypeClient}
   c.AddContact(personentity.NewContact(contactAttr))
   assert.Equal(t, c.ID, c.Contact.ObjectID)
   addrAttr := &addressentity.AddressCommonAttributes{Street: "S", Neighborhood: "N", Cep: "00000"}
   addr := addressentity.NewAddress(addrAttr)
   c.AddAddress(addr)
   assert.Equal(t, c.ID, c.Address.ObjectID)
}