package companycategoryentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
)

type CategoryToAdvertising struct {
	bun.BaseModel     `bun:"table:category_to_advertising"`
	CompanyCategoryID uuid.UUID                      `bun:"type:uuid,pk"`
	CompanyCategory   *CompanyCategory               `bun:"rel:belongs-to,join:company_category_id=id"`
	AdvertisingID     uuid.UUID                      `bun:"type:uuid,pk"`
	Advertising       *advertisingentity.Advertising `bun:"rel:belongs-to,join:advertising_id=id"`
	DeletedAt         time.Time                      `bun:",soft_delete,nullzero"`
}
