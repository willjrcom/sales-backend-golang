package advertisingentity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAdvertising(t *testing.T) {
	aca := AdvertisingCommonAttributes{Title: "T", Description: "D", Link: "L", Contact: "C", CoverImagePath: "C", Images: []string{"I"}}
	a := NewAdvertising(aca)
	assert.NotEqual(t, uuid.Nil, a.ID)
	assert.Equal(t, "T", a.Title)
	assert.Equal(t, "D", a.Description)
	assert.Equal(t, "L", a.Link)
	assert.Equal(t, "C", a.Contact)
	assert.Equal(t, "C", a.CoverImagePath)
	assert.Equal(t, []string{"I"}, a.Images)
}
