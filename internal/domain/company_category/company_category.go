package companycategoryentity

import (
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
)

type CompanyCategory struct {
	entity.Entity
	CompanyCategoryCommonAttributes
}

type CompanyCategoryCommonAttributes struct {
	Name                       string
	ImagePath                  string
	CompanyCategorySponsor     []sponsorentity.Sponsor
	CompanyCategoryAdvertising []advertisingentity.Advertising
}

type PatchCompanyCategory struct {
	Name                       *string                         `json:"name"`
	ImagePath                  *string                         `json:"image_path"`
	CompanyCategorySponsor     []sponsorentity.Sponsor         `json:"company_category_to_sponsor"`
	CompanyCategoryAdvertising []advertisingentity.Advertising `json:"company_category_to_advertising"`
}

func NewCategory(companyCategoryCommonAttributes CompanyCategoryCommonAttributes) *CompanyCategory {
	return &CompanyCategory{
		Entity:                          entity.NewEntity(),
		CompanyCategoryCommonAttributes: companyCategoryCommonAttributes,
	}
}
