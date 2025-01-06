package companydto

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyNameDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Cnpj string    `json:"cnpj"`
}

func (c *CompanyNameDTO) FromDomain(company *companyentity.Company) {
	*c = CompanyNameDTO{
		ID:   company.ID,
		Name: company.TradeName,
		Cnpj: company.Cnpj,
	}
}
