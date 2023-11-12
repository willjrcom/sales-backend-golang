package itementity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

type GroupItem struct {
	entity.Entity
	GroupDetails
	bun.BaseModel `bun:"table:group_items"`
	Items         []Item    `bun:"rel:has-many,join:id=group_item_id"`
	OrderID       uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
}

type GroupDetails struct {
	Size       float64                 `bun:"size"`
	Status     StatusItem              `bun:"status,notnull"`
	LaunchedAt *time.Time              `bun:"launch"`
	CategoryID uuid.UUID               `bun:"column:category_id,type:uuid,notnull"`
	Category   *productentity.Category `bun:"rel:belongs-to"`
}

func (i *GroupItem) CalculateQuantity() float64 {
	total := 0.0
	for _, item := range i.Items {
		total += item.Quantity
	}
	return total
}

func (i *GroupItem) CalculatePrice() float64 {
	total := 0.0
	for _, item := range i.Items {
		total += item.Price * item.Quantity
	}
	return total
}
