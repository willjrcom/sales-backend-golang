package processdto

import (
	"time"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

type OrderProcessDTO struct {
	ID                uuid.UUID                        `json:"id,omitempty"`
	EmployeeID        *uuid.UUID                       `json:"employee_id,omitempty"`
	GroupItemID       uuid.UUID                        `json:"group_item_id,omitempty"`
	GroupItem         *groupitemdto.GroupItemDTO       `json:"group_item,omitempty"`
	ProcessRuleID     uuid.UUID                        `json:"process_rule_id,omitempty"`
	Status            orderprocessentity.StatusProcess `json:"status,omitempty"`
	Products          []productcategorydto.ProductDTO  `json:"products,omitempty"`
	Queue             *orderqueuedto.QueueDTO          `json:"queue,omitempty"`
	StartedAt         *time.Time                       `json:"started_at,omitempty"`
	PausedAt          *time.Time                       `json:"paused_at,omitempty"`
	ContinuedAt       *time.Time                       `json:"continued_at,omitempty"`
	FinishedAt        *time.Time                       `json:"finished_at,omitempty"`
	CanceledAt        *time.Time                       `json:"canceled_at,omitempty"`
	CanceledReason    *string                          `json:"canceled_reason,omitempty"`
	Duration          time.Duration                    `json:"duration,omitempty"`
	DurationFormatted string                           `json:"duration_formatted,omitempty"`
	TotalPaused       int8                             `json:"total_paused,omitempty"`
}

func (s *OrderProcessDTO) FromDomain(orderProcess *orderprocessentity.OrderProcess) {
	if orderProcess == nil {
		return
	}
	*s = OrderProcessDTO{
		ID:                orderProcess.ID,
		EmployeeID:        orderProcess.EmployeeID,
		GroupItemID:       orderProcess.GroupItemID,
		GroupItem:         &groupitemdto.GroupItemDTO{},
		ProcessRuleID:     orderProcess.ProcessRuleID,
		Status:            orderProcess.Status,
		StartedAt:         orderProcess.StartedAt,
		PausedAt:          orderProcess.PausedAt,
		ContinuedAt:       orderProcess.ContinuedAt,
		FinishedAt:        orderProcess.FinishedAt,
		CanceledAt:        orderProcess.CanceledAt,
		CanceledReason:    orderProcess.CanceledReason,
		Duration:          orderProcess.Duration,
		DurationFormatted: orderProcess.Duration.String(),
		TotalPaused:       orderProcess.TotalPaused,
		Products:          []productcategorydto.ProductDTO{},
		Queue:             &orderqueuedto.QueueDTO{},
	}

	s.GroupItem.FromDomain(orderProcess.GroupItem)
	s.Queue.FromDomain(orderProcess.Queue)

	for _, product := range orderProcess.Products {
		p := &productcategorydto.ProductDTO{}
		p.FromDomain(&product)
		s.Products = append(s.Products, *p)
	}
}
