package companydto

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanySchemaDTO struct {
	ID     uuid.UUID `json:"company_id"`
	Schema *string   `json:"schema"`
}

func (c *CompanySchemaDTO) FromDomain(company *companyentity.Company) {
	if company == nil {
		return
	}
	*c = CompanySchemaDTO{
		ID:     company.ID,
		Schema: &company.SchemaName,
	}
}
