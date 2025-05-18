package orderprocessentity

import (
   "testing"
   "time"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewOrderQueue(t *testing.T) {
   gid := uuid.New()
   now := time.Now()
   q, err := NewOrderQueue(gid, now)
   assert.NoError(t, err)
   assert.NotEqual(t, uuid.Nil, q.ID)
   assert.Equal(t, gid, q.GroupItemID)
   assert.Equal(t, now, q.JoinedAt)
   assert.Equal(t, "0s", q.DurationFormatted)
}

func TestFinishQueue(t *testing.T) {
   gid := uuid.New()
   now := time.Now()
   q, _ := NewOrderQueue(gid, now)
   later := now.Add(2 * time.Second)
   q.Finish(uuid.New(), later)
   assert.NotNil(t, q.LeftAt)
   assert.Equal(t, later, *q.LeftAt)
   assert.Equal(t, q.LeftAt.Sub(q.JoinedAt), q.Duration)
   assert.Equal(t, q.Duration.String(), q.DurationFormatted)
}