package personentity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractNumber(t *testing.T) {
	n1, n2, err := ValidateAndExtractContact("(11) 96384-9111")

	assert.Nil(t, err)
	assert.Equal(t, "11", n1)
	assert.Equal(t, "963849111", n2)

	n1, n2, err = ValidateAndExtractContact("(11) 963849111")

	assert.Nil(t, err)
	assert.Equal(t, "11", n1)
	assert.Equal(t, "963849111", n2)

	n1, n2, err = ValidateAndExtractContact("(11) 96384 9111")

	assert.Nil(t, err)
	assert.Equal(t, "11", n1)
	assert.Equal(t, "963849111", n2)

	// Err
	n1, n2, err = ValidateAndExtractContact("11 96384-9111")

	assert.NotNil(t, err)
	assert.Equal(t, "", n1)
	assert.Equal(t, "", n2)

	n1, n2, err = ValidateAndExtractContact("11 96384-91111")

	assert.NotNil(t, err)
	assert.Equal(t, "", n1)
	assert.Equal(t, "", n2)

	n1, n2, err = ValidateAndExtractContact("11 963845-9111")

	assert.NotNil(t, err)
	assert.Equal(t, "", n1)
	assert.Equal(t, "", n2)

	n1, n2, err = ValidateAndExtractContact("")

	assert.NotNil(t, err)
	assert.Equal(t, "", n1)
	assert.Equal(t, "", n2)
}
