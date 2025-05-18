package employeeentity

import (
   "testing"
   "time"
   "github.com/shopspring/decimal"
   "github.com/google/uuid"
   "github.com/stretchr/testify/assert"
)

func TestNewPaymentEmployee(t *testing.T) {
   payDate := time.Now().UTC()
   p := NewPaymentEmployee(uuid.New(), decimal.NewFromFloat(100.0), StatusCompleted, MethodCash, payDate, "notes")
   assert.NotEqual(t, uuid.Nil, p.ID)
   assert.Equal(t, decimal.NewFromFloat(100.0), p.Amount)
   assert.Equal(t, StatusCompleted, p.Status)
   assert.Equal(t, MethodCash, p.Method)
   assert.Equal(t, payDate, p.PayDate)
   assert.Equal(t, "notes", p.Notes)
}