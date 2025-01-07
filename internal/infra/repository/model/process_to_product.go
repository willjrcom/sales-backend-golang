package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type OrderProcessToProductToGroupItem struct {
	bun.BaseModel `bun:"table:process_to_product_to_group_item"`
	ProcessID     uuid.UUID     `bun:"type:uuid,pk"`
	Process       *OrderProcess `bun:"rel:belongs-to,join:process_id=id"`
	ProductID     uuid.UUID     `bun:"type:uuid,pk"`
	Product       *Product      `bun:"rel:belongs-to,join:product_id=id"`
	GroupItemID   uuid.UUID     `bun:"type:uuid,pk"`
	GroupItem     *GroupItem    `bun:"rel:belongs-to,join:group_item_id=id"`
}

func (op *OrderProcessToProductToGroupItem) FromDomain(orderProcessToProductToGroupItem *orderprocessentity.OrderProcessToProductToGroupItem) {
	if orderProcessToProductToGroupItem == nil {
		return
	}
	*op = OrderProcessToProductToGroupItem{
		ProcessID:   orderProcessToProductToGroupItem.ProcessID,
		Process:     &OrderProcess{},
		ProductID:   orderProcessToProductToGroupItem.ProductID,
		Product:     &Product{},
		GroupItemID: orderProcessToProductToGroupItem.GroupItemID,
		GroupItem:   &GroupItem{},
	}

	op.Process.FromDomain(orderProcessToProductToGroupItem.Process)
	op.Product.FromDomain(orderProcessToProductToGroupItem.Product)
	op.GroupItem.FromDomain(orderProcessToProductToGroupItem.GroupItem)
}

func (op *OrderProcessToProductToGroupItem) ToDomain() *orderprocessentity.OrderProcessToProductToGroupItem {
	if op == nil {
		return nil
	}
	return &orderprocessentity.OrderProcessToProductToGroupItem{
		ProcessID:   op.ProcessID,
		Process:     op.Process.ToDomain(),
		ProductID:   op.ProductID,
		Product:     op.Product.ToDomain(),
		GroupItemID: op.GroupItemID,
		GroupItem:   op.GroupItem.ToDomain(),
	}
}
