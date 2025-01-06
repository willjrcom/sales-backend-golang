package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type GroupItem struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_group_items"`
	GroupCommonAttributes
}

type GroupCommonAttributes struct {
	GroupDetails
	GroupItemTimeLogs
	Items   []Item    `bun:"rel:has-many,join:id=group_item_id"`
	OrderID uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
}

type GroupDetails struct {
	Size             string           `bun:"size,notnull"`
	Status           StatusGroupItem  `bun:"status,notnull"`
	TotalPrice       float64          `bun:"total_price"`
	Quantity         float64          `bun:"quantity"`
	NeedPrint        bool             `bun:"need_print"`
	UseProcessRule   bool             `bun:"use_process_rule"`
	Observation      string           `bun:"observation"`
	CategoryID       uuid.UUID        `bun:"column:category_id,type:uuid,notnull"`
	Category         *ProductCategory `bun:"rel:belongs-to"`
	ComplementItemID *uuid.UUID       `bun:"column:complement_item_id,type:uuid"`
	ComplementItem   *Item            `bun:"rel:belongs-to"`
}

type GroupItemTimeLogs struct {
	StartAt    *time.Time `bun:"start_at"`
	PendingAt  *time.Time `bun:"pending_at"`
	StartedAt  *time.Time `bun:"started_at"`
	ReadyAt    *time.Time `bun:"ready_at"`
	CanceledAt *time.Time `bun:"canceled_at"`
}

func (i *GroupItem) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	if _, ok := query.GetModel().Value().(*GroupItem); ok {
		//model.CalculateTotalPrice()
	}
	return nil
}
