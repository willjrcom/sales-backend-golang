package companyentity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	cnpjservice "github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
)

func TestNewAndUpdateCompany(t *testing.T) {
	data := &cnpjservice.Cnpj{
		Cnpj: "123", BusinessName: "B", TradeName: "T", Street: "St",
		Number: "1", Neighborhood: "Nb", City: "C", UF: "U", Cep: "P",
	}
	c := NewCompany(data)
	assert.NotEqual(t, uuid.Nil, c.ID)
	assert.Contains(t, c.SchemaName, "company_")
	// update values
	newData := *data
	newData.TradeName = "TT"
	c.UpdateCompany(&newData)
	assert.Equal(t, "TT", c.TradeName)
	assert.Equal(t, "TT", c.TradeName)
}
