package companycategoryentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
)

type CategoryToSponsor struct {
	bun.BaseModel     `bun:"table:category_to_sponsor"`
	CompanyCategoryID uuid.UUID              `bun:"type:uuid,pk"`
	CompanyCategory   *CompanyCategory       `bun:"rel:belongs-to,join:company_category_id=id"`
	SponsorID         uuid.UUID              `bun:"type:uuid,pk"`
	Sponsor           *sponsorentity.Sponsor `bun:"rel:belongs-to,join:sponsor_id=id"`
	DeletedAt         time.Time              `bun:",soft_delete,nullzero"`
}
