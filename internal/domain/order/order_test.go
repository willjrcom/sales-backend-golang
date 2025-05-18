package orderentity

import (
   "testing"

   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewDefaultOrderAndStaging(t *testing.T) {
   shiftID := uuid.New()
   orderNum := 5
   attendant := uuid.New()
   o := NewDefaultOrder(shiftID, orderNum, &attendant)
   assert.NotEqual(t, uuid.Nil, o.ID)
   assert.Equal(t, OrderStatusStaging, o.Status)
   assert.Equal(t, orderNum, o.OrderNumber)
   assert.Equal(t, shiftID, o.OrderDetail.ShiftID)
}

func TestPendingOrder_NoItems_Error(t *testing.T) {
   o := NewDefaultOrder(uuid.New(), 1, nil)
   err := o.PendingOrder()
   assert.Equal(t, ErrOrderWithoutItems, err)
}