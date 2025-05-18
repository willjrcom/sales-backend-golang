package schemaentity

import (
   "testing"

   "github.com/stretchr/testify/assert"
)

func TestSchemaConstants(t *testing.T) {
   assert.Equal(t, "public", DEFAULT_SCHEMA)
   assert.Equal(t, "lost", LOST_SCHEMA)
}