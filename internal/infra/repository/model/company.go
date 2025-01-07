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

func (c *Company) FromDomain(company *companyentity.Company) {
	if company == nil {
		return
	}
	*c = Company{
		Entity: entitymodel.FromDomain(company.Entity),
		CompanyCommonAttributes: CompanyCommonAttributes{
			SchemaName:   company.SchemaName,
			BusinessName: company.BusinessName,
			TradeName:    company.TradeName,
			Cnpj:         company.Cnpj,
			Email:        company.Email,
			Contacts:     company.Contacts,
		},
	}

	c.Address.FromDomain(company.Address)
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
