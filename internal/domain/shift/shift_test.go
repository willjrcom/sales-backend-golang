package shiftentity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewShift(t *testing.T) {
	s := NewShift(decimal.NewFromFloat(10.5))
	assert.NotEqual(t, uuid.Nil, s.ID)
	assert.Equal(t, 0, s.CurrentOrderNumber)
	assert.NotNil(t, s.OpenedAt)
	assert.WithinDuration(t, time.Now().UTC(), *s.OpenedAt, time.Second)
}

func TestCloseShiftAndIsClosed(t *testing.T) {
	s := NewShift(decimal.Zero)
	assert.False(t, s.IsClosed())
	s.CloseShift(decimal.NewFromFloat(100))
	assert.True(t, s.IsClosed())
	assert.Equal(t, decimal.NewFromFloat(100), *s.EndChange)
	assert.NotNil(t, s.ClosedAt)
}

func TestIncrementCurrentOrder(t *testing.T) {
	s := NewShift(decimal.Zero)
	s.IncrementCurrentOrder()
	assert.Equal(t, 1, s.CurrentOrderNumber)
}

func TestAddRedeem(t *testing.T) {
	s := NewShift(decimal.Zero)
	r := &Redeem{Name: "r", Value: decimal.NewFromFloat(2.5)}
	s.AddRedeem(r)
	assert.Len(t, s.Redeems, 1)
	assert.Equal(t, "r", s.Redeems[0].Name)
	assert.Equal(t, decimal.NewFromFloat(2.5), s.Redeems[0].Value)
}
