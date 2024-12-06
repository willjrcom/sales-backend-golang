package companycategoryentity

import (
	"time"

	"github.com/uptrace/bun"
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
)

type CompanyCategory struct {
	entity.Entity
	bun.BaseModel `bun:"table:company_categories"`
	CompanyCategoryCommonAttributes
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type CompanyCategoryCommonAttributes struct {
	Name                       string                          `bun:"name,unique,notnull" json:"name"`
	ImagePath                  string                          `bun:"image_path" json:"image_path"`
	CompanyCategorySponsor     []sponsorentity.Sponsor         `bun:"m2m:company_category_to_sponsor,join:CompanyCategory=Sponsor" json:"company_category_to_sponsor,omitempty"`
	CompanyCategoryAdvertising []advertisingentity.Advertising `bun:"m2m:company_category_to_advertising,join:CompanyCategory=Sponsor" json:"company_category_to_advertising,omitempty"`
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
