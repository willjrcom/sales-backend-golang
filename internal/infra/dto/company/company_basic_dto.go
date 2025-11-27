package companydto

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

// CompanyBasicDTO exposes only the minimal fields required for analytics/overview use cases.
type CompanyBasicDTO struct {
	ID           uuid.UUID `json:"id"`
	BusinessName string    `json:"business_name"`
	TradeName    string    `json:"trade_name"`
	Email        string    `json:"email"`
	Cnpj         string    `json:"cnpj"`
	SchemaName   string    `json:"schema_name"`
}

func (c *CompanyBasicDTO) FromDomain(company *companyentity.Company) {
	if company == nil {
		return
	}

	*c = CompanyBasicDTO{
		ID:           company.ID,
		BusinessName: company.BusinessName,
		TradeName:    company.TradeName,
		Email:        company.Email,
		Cnpj:         company.Cnpj,
		SchemaName:   company.SchemaName,
	}
}
