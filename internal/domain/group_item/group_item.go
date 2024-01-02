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
	bun.BaseModel `bun:"table:group_items"`
	GroupCommonAttributes
}

type GroupCommonAttributes struct {
	GroupDetails
	Items   []itementity.Item `bun:"rel:has-many,join:id=group_item_id" json:"items"`
	OrderID uuid.UUID         `bun:"column:order_id,type:uuid,notnull" json:"order_id"`
}

type GroupDetails struct {
	Size       string                  `bun:"size,notnull" json:"size"`
	Status     StatusItem              `bun:"status,notnull" json:"status"`
	Price      float64                 `bun:"price" json:"price"`
	Quantity   float64                 `bun:"quantity" json:"quantity"`
	LaunchedAt *time.Time              `bun:"launch" json:"launch"`
	CategoryID uuid.UUID               `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Category   *productentity.Category `bun:"rel:belongs-to" json:"category"`
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
