package employeeentity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewPaymentEmployee(t *testing.T) {
	paymentDate := time.Now().UTC()
	p := NewPaymentEmployee(uuid.New(), decimal.NewFromFloat(100.0), StatusCompleted, "", paymentDate, "notes")
	assert.NotEqual(t, uuid.Nil, p.ID)
	assert.Equal(t, decimal.NewFromFloat(100.0), p.Amount)
	assert.Equal(t, StatusCompleted, p.Status)
	assert.Equal(t, paymentDate, p.PaymentDate)
	assert.Equal(t, "notes", p.Notes)
}
