package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderProcessToProductToGroupItem struct {
	bun.BaseModel `bun:"table:process_to_product_to_group_item"`
	ProcessID     uuid.UUID              `bun:"type:uuid,pk"`
	Process       *OrderProcess          `bun:"rel:belongs-to,join:process_id=id"`
	ProductID     uuid.UUID              `bun:"type:uuid,pk"`
	GroupItemID   uuid.UUID              `bun:"type:uuid,pk"`
	GroupItem     *orderentity.GroupItem `bun:"rel:belongs-to,join:group_item_id=id"`
}
