package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type GroupItem struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_group_items"`
	GroupCommonAttributes
	GroupItemTimeLogs
}

type GroupCommonAttributes struct {
	GroupDetails
	Items   []Item    `bun:"rel:has-many,join:id=group_item_id"`
	OrderID uuid.UUID `bun:"column:order_id,type:uuid,notnull"`
}

type GroupDetails struct {
	Size             string           `bun:"size,notnull"`
	Status           string           `bun:"status,notnull"`
	TotalPrice       decimal.Decimal  `bun:"total_price,type:decimal(10,2)"`
	Quantity         float64          `bun:"quantity"`
	NeedPrint        bool             `bun:"need_print"`
	PrinterName      string           `bun:"printer_name"`
	UseProcessRule   bool             `bun:"use_process_rule"`
	Observation      string           `bun:"observation"`
	CategoryID       uuid.UUID        `bun:"column:category_id,type:uuid,notnull"`
	Category         *ProductCategory `bun:"rel:belongs-to,join:category_id=id"`
	ComplementItemID *uuid.UUID       `bun:"column:complement_item_id,type:uuid"`
	ComplementItem   *Item            `bun:"rel:belongs-to,join:complement_item_id=id"`
}

type GroupItemTimeLogs struct {
	StartAt    *time.Time `bun:"start_at"`
	PendingAt  *time.Time `bun:"pending_at"`
	StartedAt  *time.Time `bun:"started_at"`
	ReadyAt    *time.Time `bun:"ready_at"`
	CanceledAt *time.Time `bun:"canceled_at"`
}

func (g *GroupItem) FromDomain(groupItem *orderentity.GroupItem) {
	if groupItem == nil {
		return
	}
	*g = GroupItem{
		Entity: entitymodel.FromDomain(groupItem.Entity),
		GroupCommonAttributes: GroupCommonAttributes{
			OrderID: groupItem.OrderID,
			Items:   []Item{},
			GroupDetails: GroupDetails{
				Size:             groupItem.Size,
				Status:           string(groupItem.Status),
				TotalPrice:       groupItem.TotalPrice,
				Quantity:         groupItem.Quantity,
				NeedPrint:        groupItem.NeedPrint,
				PrinterName:      groupItem.PrinterName,
				UseProcessRule:   groupItem.UseProcessRule,
				Observation:      groupItem.Observation,
				CategoryID:       groupItem.CategoryID,
				Category:         &ProductCategory{},
				ComplementItemID: groupItem.ComplementItemID,
				ComplementItem:   &Item{},
			},
		},
		GroupItemTimeLogs: GroupItemTimeLogs{
			StartAt:    groupItem.StartAt,
			PendingAt:  groupItem.PendingAt,
			StartedAt:  groupItem.StartedAt,
			ReadyAt:    groupItem.ReadyAt,
			CanceledAt: groupItem.CanceledAt,
		},
	}

	g.Category.FromDomain(groupItem.Category)
	g.ComplementItem.FromDomain(groupItem.ComplementItem)

	for _, item := range groupItem.Items {
		i := Item{}
		i.FromDomain(&item)
		g.Items = append(g.Items, i)
	}
}

func (g *GroupItem) ToDomain() *orderentity.GroupItem {
	if g == nil {
		return nil
	}
	groupItem := &orderentity.GroupItem{
		Entity: g.Entity.ToDomain(),
		GroupCommonAttributes: orderentity.GroupCommonAttributes{
			OrderID: g.OrderID,
			Items:   []orderentity.Item{},
			GroupDetails: orderentity.GroupDetails{
				Size:             g.Size,
				Status:           orderentity.StatusGroupItem(g.Status),
				TotalPrice:       g.TotalPrice,
				Quantity:         g.Quantity,
				NeedPrint:        g.NeedPrint,
				PrinterName:      g.PrinterName,
				UseProcessRule:   g.UseProcessRule,
				Observation:      g.Observation,
				CategoryID:       g.CategoryID,
				Category:         g.Category.ToDomain(),
				ComplementItemID: g.ComplementItemID,
				ComplementItem:   g.ComplementItem.ToDomain(),
			},
		},
		GroupItemTimeLogs: orderentity.GroupItemTimeLogs{
			StartAt:    g.StartAt,
			PendingAt:  g.PendingAt,
			StartedAt:  g.StartedAt,
			ReadyAt:    g.ReadyAt,
			CanceledAt: g.CanceledAt,
		},
	}

	for _, item := range g.Items {
		groupItem.Items = append(groupItem.Items, *item.ToDomain())
	}

	return groupItem
}
