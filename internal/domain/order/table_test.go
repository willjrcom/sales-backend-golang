package orderentity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockUnlockTable(t *testing.T) {
	tble := &Table{TableCommonAttributes: TableCommonAttributes{Name: "T", IsAvailable: true}}
	tble.LockTable()
	assert.False(t, tble.IsAvailable)
	tble.UnlockTable()
	assert.True(t, tble.IsAvailable)
}
