package groupitemdto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

type GroupItemDTO struct {
	GroupItemTimeLogsDTO
	ID               uuid.UUID                       `json:"id"`
	Size             string                          `json:"size"`
	Status           orderentity.StatusGroupItem     `json:"status"`
	TotalPrice       decimal.Decimal                 `json:"total_price"`
	Quantity         float64                         `json:"quantity"`
	NeedPrint        bool                            `json:"need_print"`
	PrinterName      string                          `json:"printer_name"`
	UseProcessRule   bool                            `json:"use_process_rule"`
	Observation      string                          `json:"observation"`
	CategoryID       uuid.UUID                       `json:"category_id"`
	Category         *productcategorydto.CategoryDTO `json:"category"`
	ComplementItemID *uuid.UUID                      `json:"complement_item_id"`
	ComplementItem   *itemdto.ItemDTO                `json:"complement_item"`
	Items            []itemdto.ItemDTO               `json:"items"`
	OrderID          uuid.UUID                       `json:"order_id"`
}

type GroupItemTimeLogsDTO struct {
	StartAt    *time.Time `json:"start_at"`
	PendingAt  *time.Time `json:"pending_at"`
	StartedAt  *time.Time `json:"started_at"`
	ReadyAt    *time.Time `json:"ready_at"`
	CanceledAt *time.Time `json:"canceled_at"`
}

func (i *GroupItemDTO) FromDomain(groupItem *orderentity.GroupItem) {
	if groupItem == nil {
		return
	}
	*i = GroupItemDTO{
		ID:               groupItem.ID,
		Size:             groupItem.Size,
		Status:           groupItem.Status,
		TotalPrice:       groupItem.TotalPrice,
		Quantity:         groupItem.Quantity,
		NeedPrint:        groupItem.NeedPrint,
		PrinterName:      groupItem.PrinterName,
		UseProcessRule:   groupItem.UseProcessRule,
		Observation:      groupItem.Observation,
		CategoryID:       groupItem.CategoryID,
		Category:         &productcategorydto.CategoryDTO{},
		ComplementItemID: groupItem.ComplementItemID,
		ComplementItem:   &itemdto.ItemDTO{},
		Items:            []itemdto.ItemDTO{},
		OrderID:          groupItem.OrderID,
		GroupItemTimeLogsDTO: GroupItemTimeLogsDTO{
			StartAt:    groupItem.StartAt,
			PendingAt:  groupItem.PendingAt,
			StartedAt:  groupItem.StartedAt,
			ReadyAt:    groupItem.ReadyAt,
			CanceledAt: groupItem.CanceledAt,
		},
	}

	i.Category.FromDomain(groupItem.Category)
	i.ComplementItem.FromDomain(groupItem.ComplementItem)

	for _, item := range groupItem.Items {
		itemDTO := itemdto.ItemDTO{}
		itemDTO.FromDomain(&item)
		i.Items = append(i.Items, itemDTO)
	}

	if groupItem.Category == nil {
		i.Category = nil
	}
	if groupItem.ComplementItem == nil {
		i.ComplementItem = nil
	}
	if len(groupItem.Items) == 0 {
		i.Items = nil
	}
}
