package companycategoryentity

import (
	"github.com/google/uuid"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
)

type CategoryToSponsor struct {
	CompanyCategoryID uuid.UUID
	CompanyCategory   *CompanyCategory
	SponsorID         uuid.UUID
	Sponsor           *sponsorentity.Sponsor
}
