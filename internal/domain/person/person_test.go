package personentity

import (
   "testing"
   "time"

   "github.com/stretchr/testify/assert"
)

func TestNewPerson(t *testing.T) {
   birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
   common := &PersonCommonAttributes{
       Name:     "Name",
       Email:    "e@example.com",
       Cpf:      "123",
       Birthday: &birthday,
   }
   p := NewPerson(common)
   assert.Equal(t, "Name", p.Name)
   assert.Equal(t, "e@example.com", p.Email)
   assert.Equal(t, &birthday, p.Birthday)
}