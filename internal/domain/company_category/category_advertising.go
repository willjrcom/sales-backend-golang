package companycategoryentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
)

type CategoryToAdvertising struct {
	bun.BaseModel     `bun:"table:category_to_advertising"`
	CategoryCompanyID uuid.UUID                      `bun:"type:uuid,pk"`
	CategoryCompany   *CompanyCategory               `bun:"rel:belongs-to,join:category_company_id=id"`
	AdvertisingID     uuid.UUID                      `bun:"type:uuid,pk"`
	Advertising       *advertisingentity.Advertising `bun:"rel:belongs-to,join:advertising_id=id"`
}
