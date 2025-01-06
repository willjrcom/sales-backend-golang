package model

import (
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanyWithUsers struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:companies"`
	CompanyCommonAttributes
	Users []User `bun:"m2m:company_to_users,join:CompanyWithUsers=User"`
}

func (c *CompanyWithUsers) FromDomain(model *companyentity.CompanyWithUsers) {
	*c = CompanyWithUsers{
		Entity: entitymodel.FromDomain(model.Entity),
		CompanyCommonAttributes: CompanyCommonAttributes{
			SchemaName:   model.SchemaName,
			BusinessName: model.BusinessName,
			TradeName:    model.TradeName,
			Cnpj:         model.Cnpj,
			Email:        model.Email,
			Contacts:     model.Contacts,
		},
		Users: []User{},
	}

	for _, user := range model.Users {
		u := User{}
		u.FromDomain(&user)
		c.Users = append(c.Users, u)
	}
}

func (c *CompanyWithUsers) ToDomain() *companyentity.CompanyWithUsers {
	if c == nil {
		return nil
	}
	company := &companyentity.CompanyWithUsers{
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
		Users: []companyentity.User{},
	}

	for _, user := range c.Users {
		company.Users = append(company.Users, *user.ToDomain())
	}

	return company
}
