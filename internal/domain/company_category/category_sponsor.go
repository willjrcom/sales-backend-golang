package companycategoryentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
)

type CategoryToSponsor struct {
	bun.BaseModel     `bun:"table:category_to_sponsor"`
	CategoryCompanyID uuid.UUID              `bun:"type:uuid,pk"`
	CategoryCompany   *CompanyCategory       `bun:"rel:belongs-to,join:category_company_id=id"`
	SponsorID         uuid.UUID              `bun:"type:uuid,pk"`
	Sponsor           *sponsorentity.Sponsor `bun:"rel:belongs-to,join:sponsor_id=id"`
}
