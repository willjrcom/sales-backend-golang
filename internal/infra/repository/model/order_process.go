package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type OrderProcess struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_processes,alias:process"`
	OrderProcessTimeLogs
	OrderProcessCommonAttributes
}

type OrderProcessCommonAttributes struct {
	EmployeeID    *uuid.UUID  `bun:"employee_id,type:uuid"`
	GroupItemID   uuid.UUID   `bun:"group_item_id,type:uuid,notnull"`
	GroupItem     *GroupItem  `bun:"rel:belongs-to,join:group_item_id=id"`
	ProcessRuleID uuid.UUID   `bun:"process_rule_id,type:uuid,notnull"`
	Status        string      `bun:"status,notnull"`
	Products      []Product   `bun:"m2m:process_to_product_to_group_item,join:Process=Product"`
	Queue         *OrderQueue `bun:"rel:has-one,join:group_item_id=group_item_id,join:process_rule_id=process_rule_id"`
}

type OrderProcessTimeLogs struct {
	StartedAt         *time.Time    `bun:"started_at"`
	PausedAt          *time.Time    `bun:"paused_at"`
	ContinuedAt       *time.Time    `bun:"continued_at"`
	FinishedAt        *time.Time    `bun:"finished_at"`
	CanceledAt        *time.Time    `bun:"canceled_at"`
	CanceledReason    *string       `bun:"canceled_reason"`
	Duration          time.Duration `bun:"duration"`
	DurationFormatted string        `bun:"duration_formatted"`
	TotalPaused       int8          `bun:"total_paused"`
}

func (op *OrderProcess) FromDomain(orderProcess *orderprocessentity.OrderProcess) {
	if orderProcess == nil {
		return
	}
	*op = OrderProcess{
		Entity: entitymodel.Entity{
			ID:        orderProcess.ID,
			CreatedAt: orderProcess.CreatedAt,
			UpdatedAt: orderProcess.UpdatedAt,
		},
		OrderProcessTimeLogs: OrderProcessTimeLogs{
			StartedAt:         orderProcess.StartedAt,
			PausedAt:          orderProcess.PausedAt,
			ContinuedAt:       orderProcess.ContinuedAt,
			FinishedAt:        orderProcess.FinishedAt,
			CanceledAt:        orderProcess.CanceledAt,
			CanceledReason:    orderProcess.CanceledReason,
			Duration:          orderProcess.Duration,
			DurationFormatted: orderProcess.Duration.String(),
			TotalPaused:       orderProcess.TotalPaused,
		},
		OrderProcessCommonAttributes: OrderProcessCommonAttributes{
			EmployeeID:    orderProcess.EmployeeID,
			GroupItemID:   orderProcess.GroupItemID,
			GroupItem:     &GroupItem{},
			ProcessRuleID: orderProcess.ProcessRuleID,
			Status:        string(orderProcess.Status),
			Products:      []Product{},
			Queue:         &OrderQueue{},
		},
	}

	op.GroupItem.FromDomain(orderProcess.GroupItem)
	op.Queue.FromDomain(orderProcess.Queue)

	for _, product := range orderProcess.Products {
		p := Product{}
		p.FromDomain(&product)
		op.Products = append(op.Products, p)
	}
}

func (op *OrderProcess) ToDomain() *orderprocessentity.OrderProcess {
	if op == nil {
		return nil
	}
	orderProcess := &orderprocessentity.OrderProcess{
		Entity: op.Entity.ToDomain(),
		OrderProcessCommonAttributes: orderprocessentity.OrderProcessCommonAttributes{
			EmployeeID:    op.EmployeeID,
			GroupItemID:   op.GroupItemID,
			GroupItem:     op.GroupItem.ToDomain(),
			ProcessRuleID: op.ProcessRuleID,
			Status:        orderprocessentity.StatusProcess(op.Status),
			Products:      []productentity.Product{},
			Queue:         op.Queue.ToDomain(),
		},
		OrderProcessTimeLogs: orderprocessentity.OrderProcessTimeLogs{
			StartedAt:         op.StartedAt,
			PausedAt:          op.PausedAt,
			ContinuedAt:       op.ContinuedAt,
			FinishedAt:        op.FinishedAt,
			CanceledAt:        op.CanceledAt,
			CanceledReason:    op.CanceledReason,
			Duration:          op.Duration,
			DurationFormatted: op.Duration.String(),
			TotalPaused:       op.TotalPaused,
		},
	}

	for _, product := range op.Products {
		orderProcess.Products = append(orderProcess.Products, *product.ToDomain())
	}

	return orderProcess
}
