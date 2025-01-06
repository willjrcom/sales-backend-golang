package model

import (
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Company struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:companies"`
	CompanyCommonAttributes
}

type CompanyCommonAttributes struct {
	SchemaName   string   `bun:"schema_name,notnull"`
	BusinessName string   `bun:"business_name,notnull"`
	TradeName    string   `bun:"trade_name,notnull"`
	Cnpj         string   `bun:"cnpj,notnull"`
	Email        string   `bun:"email"`
	Contacts     []string `bun:"contacts,type:jsonb"`
	Address      *Address `bun:"rel:has-one,join:id=object_id,notnull"`
}

func (c *Company) FromDomain(model *companyentity.Company) {
	*c = Company{
		Entity: entitymodel.FromDomain(model.Entity),
		CompanyCommonAttributes: CompanyCommonAttributes{
			SchemaName:   model.SchemaName,
			BusinessName: model.BusinessName,
			TradeName:    model.TradeName,
			Cnpj:         model.Cnpj,
			Email:        model.Email,
			Contacts:     model.Contacts,
		},
	}

	c.Address.FromDomain(model.Address)
}

func (c *Company) ToDomain() *companyentity.Company {
	if c == nil {
		return nil
	}
	return &companyentity.Company{
		Entity: c.Entity.ToDomain(),
		CompanyCommonAttributes: companyentity.CompanyCommonAttributes{
			SchemaName:   c.SchemaName,
			BusinessName: c.BusinessName,
			TradeName:    c.TradeName,
			Cnpj:         c.Cnpj,
			Email:        c.Email,
			Contacts:     c.Contacts,
			Address:      c.Address.ToDomain(),
		},
	}
}
