package personentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewContact(t *testing.T) {
   cca := &ContactCommonAttributes{Ddd: "011", Number: "7890", Type: ContactTypeEmployee}
   c := NewContact(cca)
   assert.NotEqual(t, uuid.Nil, c.ID)
   assert.Equal(t, "011", c.Ddd)
   assert.Equal(t, "7890", c.Number)
   assert.Equal(t, ContactTypeEmployee, c.Type)
}