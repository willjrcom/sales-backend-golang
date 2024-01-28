package companyentity

import (
	"strings"

	"github.com/teris-io/shortid"
	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
)

type Company struct {
	entity.Entity
	bun.BaseModel `bun:"table:companies"`
	CompanyCommonAttributes
}

type CompanyCommonAttributes struct {
	SchemaName   string                `bun:"schema_name,notnull"`
	BusinessName string                `bun:"business_name,notnull" json:"business_name"`
	TradeName    string                `bun:"trade_name,notnull" json:"trade_name"`
	Cnpj         string                `bun:"cnpj,notnull" json:"cnpj"`
	Email        string                `bun:"email" json:"email"`
	Contacts     []string              `bun:"contacts,type:jsonb" json:"contacts,omitempty"`
	Address      addressentity.Address `bun:"rel:has-one,join:id=object_id,notnull" json:"address,omitempty"`
}

func NewCompany(cnpjData *cnpj.Cnpj) *Company {
	id, _ := shortid.Generate()

	schema := "loja_" + strings.ToLower(cnpjData.TradeName) + "_" + strings.ToLower(id)
	company := &Company{
		Entity: entity.NewEntity(),
		CompanyCommonAttributes: CompanyCommonAttributes{
			BusinessName: cnpjData.BusinessName,
			TradeName:    cnpjData.TradeName,
			Cnpj:         cnpjData.Cnpj,
			SchemaName:   schema,
		},
	}

	addressCommonAttributes := &addressentity.AddressCommonAttributes{
		Street:       cnpjData.Street,
		Number:       cnpjData.Number,
		Neighborhood: cnpjData.Neighborhood,
		City:         cnpjData.City,
		State:        cnpjData.State,
		Cep:          cnpjData.Cep,
	}

	company.AddAddress(addressCommonAttributes)
	return company

}

func (c *Company) AddAddress(addressCommonAttributes *addressentity.AddressCommonAttributes) {
	addressCommonAttributes.ObjectID = c.ID
	c.Address = *addressentity.NewAddress(addressCommonAttributes)
}
