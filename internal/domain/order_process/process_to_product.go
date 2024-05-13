package orderprocessentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type OrderProcessToProductToGroupItem struct {
	bun.BaseModel `bun:"table:process_to_product_to_group_item"`
	ProcessID     uuid.UUID                  `bun:"type:uuid,pk"`
	Process       *OrderProcess              `bun:"rel:belongs-to,join:process_id=id"`
	ProductID     uuid.UUID                  `bun:"type:uuid,pk"`
	Product       *productentity.Product     `bun:"rel:belongs-to,join:product_id=id"`
	GroupItemID   uuid.UUID                  `bun:"type:uuid,pk"`
	GroupItem     *groupitementity.GroupItem `bun:"rel:belongs-to,join:group_item_id=id"`
}
