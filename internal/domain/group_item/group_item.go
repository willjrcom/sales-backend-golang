package groupitementity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type GroupItem struct {
	entity.Entity
	GroupDetails
	bun.BaseModel `bun:"table:group_items"`
	Items         []itementity.Item `bun:"rel:has-many,join:id=group_item_id"`
	OrderID       uuid.UUID         `bun:"column:order_id,type:uuid,notnull"`
}

type GroupDetails struct {
	Size       string                  `bun:"size,notnull"`
	Status     StatusItem              `bun:"status,notnull"`
	Price      float64                 `bun:"price"`
	Quantity   float64                 `bun:"quantity"`
	LaunchedAt *time.Time              `bun:"launch"`
	CategoryID uuid.UUID               `bun:"column:category_id,type:uuid,notnull"`
	Category   *productentity.Category `bun:"rel:belongs-to"`
}

func (i *GroupItem) CalculateTotalValues() {
	qtd := 0.0
	price := 0.0

	for _, item := range i.Items {
		qtd += item.Quantity
		price += item.Price * item.Quantity
	}

	i.GroupDetails.Quantity = qtd
	i.GroupDetails.Price = price
}
