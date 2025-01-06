package companydto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
)

type CompanyDTO struct {
	SchemaName   string
	BusinessName string
	TradeName    string
	Cnpj         string
	Email        string
	Contacts     []string
	Address      *addressdto.AddressDTO
}

func (o *CompanyDTO) FromModel(model *companyentity.Company) {
	*o = CompanyDTO{
		SchemaName:   model.SchemaName,
		BusinessName: model.BusinessName,
		TradeName:    model.TradeName,
		Cnpj:         model.Cnpj,
		Email:        model.Email,
		Contacts:     model.Contacts,
	}

	o.Address.FromModel(model.Address)
}
