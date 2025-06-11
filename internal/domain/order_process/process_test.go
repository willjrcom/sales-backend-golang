package orderprocessentity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderProcess(t *testing.T) {
	gid := uuid.New()
	pid := uuid.New()
	p := NewOrderProcess(gid, pid, 1, ProcessDeliveryType)
	assert.NotEqual(t, uuid.Nil, p.ID)
	assert.Equal(t, gid, p.GroupItemID)
	assert.Equal(t, pid, p.ProcessRuleID)
	assert.Equal(t, ProcessStatusPending, p.Status)
}

func TestStartAndCancelProcess(t *testing.T) {
	gid := uuid.New()
	pid := uuid.New()
	p := NewOrderProcess(gid, pid, 1, ProcessDeliveryType)
	err := p.StartProcess(uuid.Nil)
	assert.Error(t, err)
	eid := uuid.New()
	err = p.StartProcess(eid)
	assert.NoError(t, err)
	assert.Equal(t, ProcessStatusStarted, p.Status)
	reason := "reason"
	err = p.CancelProcess(nil)
	assert.Error(t, err)
	err = p.CancelProcess(&reason)
	assert.NoError(t, err)
	assert.Equal(t, ProcessStatusCanceled, p.Status)
	assert.Equal(t, &reason, p.CanceledReason)
}
