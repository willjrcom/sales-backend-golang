package entity

import (
   "testing"
   "time"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewEntity(t *testing.T) {
   e := NewEntity()
   assert.NotEqual(t, uuid.Nil, e.ID)
   now := time.Now().UTC()
   assert.WithinDuration(t, now, e.CreatedAt, time.Second)
   assert.WithinDuration(t, e.CreatedAt, e.UpdatedAt, time.Second)
   assert.Nil(t, e.DeletedAt)
}