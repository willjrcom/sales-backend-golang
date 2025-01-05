package companycategoryentity

import (
	"github.com/google/uuid"
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
)

type CategoryToAdvertising struct {
	CompanyCategoryID uuid.UUID
	CompanyCategory   *CompanyCategory
	AdvertisingID     uuid.UUID
	Advertising       *advertisingentity.Advertising
}
