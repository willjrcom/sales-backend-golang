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

func NewCategory(companyCategoryCommonAttributes CompanyCategoryCommonAttributes) *CompanyCategory {
	return &CompanyCategory{
		Entity:                          entity.NewEntity(),
		CompanyCategoryCommonAttributes: companyCategoryCommonAttributes,
	}
}
