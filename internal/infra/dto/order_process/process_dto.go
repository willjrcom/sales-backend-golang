package processdto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
)

type OrderProcessDTO struct {
	ID                uuid.UUID                           `json:"id,omitempty"`
	CreatedAt         time.Time                           `json:"created_at,omitempty"`
	OrderNumber       int                                 `json:"order_number"`
	OrderType         orderprocessentity.OrderProcessType `json:"order_type"`
	EmployeeID        *uuid.UUID                          `json:"employee_id,omitempty"`
	OrderID           uuid.UUID                           `json:"order_id,omitempty"`
	GroupItemID       uuid.UUID                           `json:"group_item_id,omitempty"`
	GroupItem         *groupitemdto.GroupItemDTO          `json:"group_item,omitempty"`
	ProcessRuleID     uuid.UUID                           `json:"process_rule_id,omitempty"`
	Status            orderprocessentity.StatusProcess    `json:"status,omitempty"`
	Products          []productcategorydto.ProductDTO     `json:"products,omitempty"`
	Queue             *orderqueuedto.QueueDTO             `json:"queue,omitempty"`
	StartedAt         *time.Time                          `json:"started_at,omitempty"`
	PausedAt          *time.Time                          `json:"paused_at,omitempty"`
	ContinuedAt       *time.Time                          `json:"continued_at,omitempty"`
	FinishedAt        *time.Time                          `json:"finished_at,omitempty"`
	CancelledAt       *time.Time                          `json:"cancelled_at,omitempty"`
	CancelledReason   *string                             `json:"cancelled_reason,omitempty"`
	Duration          time.Duration                       `json:"duration,omitempty"`
	DurationFormatted string                              `json:"duration_formatted,omitempty"`
	TotalPaused       int8                                `json:"total_paused,omitempty"`
}

func (s *OrderProcessDTO) FromDomain(orderProcess *orderprocessentity.OrderProcess) {
	if orderProcess == nil {
		return
	}

	elapsedDuration := calculateElapsedDuration(orderProcess, time.Now().UTC())

	*s = OrderProcessDTO{
		ID:                orderProcess.ID,
		CreatedAt:         orderProcess.CreatedAt,
		OrderNumber:       orderProcess.OrderNumber,
		OrderType:         orderProcess.OrderType,
		EmployeeID:        orderProcess.EmployeeID,
		OrderID:           orderProcess.OrderID,
		GroupItemID:       orderProcess.GroupItemID,
		GroupItem:         &groupitemdto.GroupItemDTO{},
		ProcessRuleID:     orderProcess.ProcessRuleID,
		Status:            orderProcess.Status,
		StartedAt:         orderProcess.StartedAt,
		PausedAt:          orderProcess.PausedAt,
		ContinuedAt:       orderProcess.ContinuedAt,
		FinishedAt:        orderProcess.FinishedAt,
		CancelledAt:       orderProcess.CancelledAt,
		CancelledReason:   orderProcess.CancelledReason,
		Duration:          elapsedDuration,
		DurationFormatted: formatDurationToClock(elapsedDuration),
		TotalPaused:       orderProcess.TotalPaused,
		Products:          []productcategorydto.ProductDTO{},
		Queue:             &orderqueuedto.QueueDTO{},
	}

	s.GroupItem.FromDomain(orderProcess.GroupItem)
	s.Queue.FromDomain(orderProcess.Queue)

	for _, product := range orderProcess.Products {
		p := productcategorydto.ProductDTO{}
		p.FromDomain(&product)
		s.Products = append(s.Products, p)
	}

	if orderProcess.GroupItem == nil {
		s.GroupItem = nil
	}
	if orderProcess.Queue == nil {
		s.Queue = nil
	}

	if len(orderProcess.Products) == 0 {
		s.Products = nil
	}
}

func calculateElapsedDuration(orderProcess *orderprocessentity.OrderProcess, now time.Time) time.Duration {
	if orderProcess == nil {
		return 0
	}

	elapsed := orderProcess.Duration

	if !isProcessRunning(orderProcess.Status) {
		return elapsed
	}

	reference := runningReferenceTime(orderProcess)
	if reference == nil || now.Before(*reference) {
		return elapsed
	}

	return elapsed + now.Sub(*reference)
}

func isProcessRunning(status orderprocessentity.StatusProcess) bool {
	return status == orderprocessentity.ProcessStatusStarted ||
		status == orderprocessentity.ProcessStatusContinued
}

func runningReferenceTime(orderProcess *orderprocessentity.OrderProcess) *time.Time {
	if orderProcess == nil {
		return nil
	}

	if orderProcess.Status == orderprocessentity.ProcessStatusContinued && orderProcess.ContinuedAt != nil {
		return orderProcess.ContinuedAt
	}

	if orderProcess.StartedAt != nil {
		return orderProcess.StartedAt
	}

	return orderProcess.ContinuedAt
}

func formatDurationToClock(duration time.Duration) string {
	if duration < 0 {
		duration = 0
	}

	totalSeconds := int(duration / time.Second)

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
