package model

import (
	"github.com/uptrace/bun"
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
	companycategoryentity "github.com/willjrcom/sales-backend-go/internal/domain/company_category"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type CompanyCategory struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:public.company_categories"`
	CompanyCategoryCommonAttributes
}

type CompanyCategoryCommonAttributes struct {
	Name           string        `bun:"name,notnull"`
	ImagePath      string        `bun:"image_path"`
	Sponsors       []Sponsor     `bun:"m2m:public.category_sponsors,join:Category=Sponsor"`
	Advertisements []Advertising `bun:"m2m:public.category_advertisements,join:Category=Advertising"`
}

func (c *CompanyCategory) FromDomain(category *companycategoryentity.CompanyCategory) {
	if category == nil {
		return
	}
	*c = CompanyCategory{
		Entity: entitymodel.FromDomain(category.Entity),
		CompanyCategoryCommonAttributes: CompanyCategoryCommonAttributes{
			Name:           category.Name,
			ImagePath:      category.ImagePath,
			Sponsors:       []Sponsor{},
			Advertisements: []Advertising{},
		},
	}

	for _, sponsor := range category.CompanyCategorySponsor {
		sponsorModel := Sponsor{}
		sponsorModel.FromDomain(&sponsor)
		c.Sponsors = append(c.Sponsors, sponsorModel)
	}

	for _, adv := range category.CompanyCategoryAdvertising {
		advModel := Advertising{}
		advModel.FromDomain(&adv)
		c.Advertisements = append(c.Advertisements, advModel)
	}
}

func (c *CompanyCategory) ToDomain() *companycategoryentity.CompanyCategory {
	if c == nil {
		return nil
	}

	sponsors := []sponsorentity.Sponsor{}
	for _, sponsor := range c.Sponsors {
		sponsors = append(sponsors, *sponsor.ToDomain())
	}

	advs := []advertisingentity.Advertising{}
	for _, adv := range c.Advertisements {
		advs = append(advs, *adv.ToDomain())
	}

	return &companycategoryentity.CompanyCategory{
		Entity: c.Entity.ToDomain(),
		CompanyCategoryCommonAttributes: companycategoryentity.CompanyCategoryCommonAttributes{
			Name:                       c.Name,
			ImagePath:                  c.ImagePath,
			CompanyCategorySponsor:     sponsors,
			CompanyCategoryAdvertising: advs,
		},
	}
}
