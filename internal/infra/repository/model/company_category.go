package model

import (
	"github.com/uptrace/bun"
	companycategoryentity "github.com/willjrcom/sales-backend-go/internal/domain/company_category"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanyCategory struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:public.company_categories"`
	CompanyCategoryCommonAttributes
}

type CompanyCategoryCommonAttributes struct {
	Name      string `bun:"name,notnull"`
	ImagePath string `bun:"image_path"`
}

func (c *CompanyCategory) FromDomain(category *companycategoryentity.CompanyCategory) {
	if category == nil {
		return
	}
	*c = CompanyCategory{
		Entity: entitymodel.FromDomain(category.Entity),
		CompanyCategoryCommonAttributes: CompanyCategoryCommonAttributes{
			Name:      category.Name,
			ImagePath: category.ImagePath,
		},
	}
}

func (c *CompanyCategory) ToDomain() *companycategoryentity.CompanyCategory {
	if c == nil {
		return nil
	}
	return &companycategoryentity.CompanyCategory{
		Entity: c.Entity.ToDomain(),
		CompanyCategoryCommonAttributes: companycategoryentity.CompanyCategoryCommonAttributes{
			Name:      c.Name,
			ImagePath: c.ImagePath,
		},
	}
}
