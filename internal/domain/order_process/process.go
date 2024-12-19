package orderprocessentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrMustBeStarted = errors.New("process must be started")
)

type OrderProcess struct {
	entity.Entity
	bun.BaseModel `bun:"table:order_processes"`
	OrderProcessTimeLogs
	OrderProcessCommonAttributes
}

type OrderProcessCommonAttributes struct {
	EmployeeID    *uuid.UUID                 `bun:"employee_id,type:uuid" json:"employee_id,omitempty"`
	GroupItemID   uuid.UUID                  `bun:"group_item_id,type:uuid,notnull" json:"group_item_id"`
	GroupItem     *groupitementity.GroupItem `bun:"rel:belongs-to,join:group_item_id=id" json:"group_item,omitempty"`
	ProcessRuleID uuid.UUID                  `bun:"process_rule_id,type:uuid,notnull" json:"process_rule_id"`
	Status        StatusProcess              `bun:"status,notnull" json:"status"`
	Products      []productentity.Product    `bun:"m2m:process_to_product_to_group_item,join:Process=Product" json:"process_to_product,omitempty"`
}

type OrderProcessTimeLogs struct {
	StartedAt         *time.Time    `bun:"started_at" json:"started_at,omitempty"`
	PausedAt          *time.Time    `bun:"paused_at" json:"paused_at,omitempty"`
	ContinuedAt       *time.Time    `bun:"continued_at" json:"continued_at,omitempty"`
	FinishedAt        *time.Time    `bun:"finished_at" json:"finished_at,omitempty"`
	Duration          time.Duration `bun:"duration" json:"duration"`
	DurationFormatted string        `bun:"duration_formatted" json:"duration_formatted"`
	TotalPaused       int8          `bun:"total_paused" json:"total_paused"`
}

func NewOrderProcess(groupItemID uuid.UUID, processRuleID uuid.UUID) *OrderProcess {
	orderProcessCommonAttributes := OrderProcessCommonAttributes{
		GroupItemID:   groupItemID,
		ProcessRuleID: processRuleID,
		Status:        ProcessStatusPending,
	}

	return &OrderProcess{
		Entity:                       entity.NewEntity(),
		OrderProcessCommonAttributes: orderProcessCommonAttributes,
		OrderProcessTimeLogs:         OrderProcessTimeLogs{},
	}
}

func (p *OrderProcess) StartProcess(employeeID uuid.UUID) error {
	if p.StartedAt != nil {
		return errors.New("process already started")
	}

	if employeeID == uuid.Nil {
		return errors.New("employee not found")
	}

	p.EmployeeID = &employeeID
	p.StartedAt = &time.Time{}
	*p.StartedAt = time.Now()
	p.Status = ProcessStatusStarted
	return nil
}

func (p *OrderProcess) FinishProcess() error {
	if p.StartedAt == nil {
		return ErrMustBeStarted
	}

	if p.PausedAt != nil {
		return errors.New("process paused, must be continue to finish")
	}

	if p.FinishedAt != nil {
		return errors.New("process already finished")
	}

	p.FinishedAt = &time.Time{}
	*p.FinishedAt = time.Now()
	p.Status = ProcessStatusFinished

	if p.ContinuedAt != nil {
		p.Duration += time.Since(*p.ContinuedAt)
		p.DurationFormatted = p.Duration.String()
		return nil
	}

	p.Duration = time.Since(*p.StartedAt)
	p.DurationFormatted = p.Duration.String()
	return nil
}

func (p *OrderProcess) PauseProcess() error {
	if p.StartedAt == nil {
		return ErrMustBeStarted
	}

	if p.PausedAt != nil {
		return errors.New("process already paused")
	}

	p.TotalPaused++

	p.PausedAt = &time.Time{}
	*p.PausedAt = time.Now()
	p.Status = ProcessStatusPaused

	if p.ContinuedAt != nil {
		p.Duration += time.Since(*p.ContinuedAt)
		p.DurationFormatted = p.Duration.String()
		p.ContinuedAt = nil
		return nil
	}

	p.Duration += time.Since(*p.StartedAt)
	p.DurationFormatted = p.Duration.String()
	return nil
}

func (p *OrderProcess) ContinueProcess() error {
	if p.StartedAt == nil {
		return ErrMustBeStarted
	}

	if p.PausedAt == nil {
		return errors.New("process must be paused")
	}

	if p.ContinuedAt != nil {
		return errors.New("process already continued")
	}

	p.ContinuedAt = &time.Time{}
	*p.ContinuedAt = time.Now()
	p.Status = ProcessStatusContinued
	p.PausedAt = nil
	return nil
}
