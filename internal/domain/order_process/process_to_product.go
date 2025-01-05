package orderprocessentity

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type OrderProcessToProductToGroupItem struct {
	ProcessID   uuid.UUID
	Process     *OrderProcess
	ProductID   uuid.UUID
	Product     *productentity.Product
	GroupItemID uuid.UUID
	GroupItem   *orderentity.GroupItem
}
