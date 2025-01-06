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

func (s *OrderProcessDTO) FromDomain(model *orderprocessentity.OrderProcess) {
	*s = OrderProcessDTO{
		ID:                model.ID,
		EmployeeID:        model.EmployeeID,
		GroupItemID:       model.GroupItemID,
		ProcessRuleID:     model.ProcessRuleID,
		Status:            model.Status,
		StartedAt:         model.StartedAt,
		PausedAt:          model.PausedAt,
		ContinuedAt:       model.ContinuedAt,
		FinishedAt:        model.FinishedAt,
		CanceledAt:        model.CanceledAt,
		CanceledReason:    model.CanceledReason,
		Duration:          model.Duration,
		DurationFormatted: model.Duration.String(),
		TotalPaused:       model.TotalPaused,
	}

	s.GroupItem.FromDomain(model.GroupItem)
	s.Queue.FromDomain(model.Queue)
}
