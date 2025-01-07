package companydto

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

type CompanyDTO struct {
	ID           uuid.UUID              `json:"id"`
	SchemaName   string                 `json:"schema_name"`
	BusinessName string                 `json:"business_name"`
	TradeName    string                 `json:"trade_name"`
	Cnpj         string                 `json:"cnpj"`
	Email        string                 `json:"email"`
	Contacts     []string               `json:"contacts"`
	Address      *addressdto.AddressDTO `json:"address"`
}

func (c *CompanyDTO) FromDomain(model *companyentity.Company) {
	*c = CompanyDTO{
		ID:           model.ID,
		SchemaName:   model.SchemaName,
		BusinessName: model.BusinessName,
		TradeName:    model.TradeName,
		Cnpj:         model.Cnpj,
		Email:        model.Email,
		Contacts:     model.Contacts,
	}

	c.Address.FromDomain(model.Address)
}
