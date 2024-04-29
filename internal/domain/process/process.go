package processentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrMustBeStarted = errors.New("process must be started")
)

type Process struct {
	entity.Entity
	bun.BaseModel `bun:"table:processes"`
	ProcessTimeLogs
	ProcessCommonAttributes
}

type ProcessCommonAttributes struct {
	EmployeeID    *uuid.UUID    `bun:"employee_id,type:uuid" json:"employee_id,omitempty"`
	ProductID     uuid.UUID     `bun:"product_id,type:uuid,notnull" json:"product_id"`
	ItemID        uuid.UUID     `bun:"item_id,type:uuid,notnull" json:"item_id"`
	ProcessRuleID uuid.UUID     `bun:"process_rule_id,type:uuid,notnull" json:"process_rule_id"`
	Status        StatusProcess `bun:"status,notnull" json:"status"`
}

type ProcessTimeLogs struct {
	StartedAt         *time.Time    `bun:"started_at" json:"started_at,omitempty"`
	PausedAt          *time.Time    `bun:"paused_at" json:"paused_at,omitempty"`
	ContinuedAt       *time.Time    `bun:"continued_at" json:"continued_at,omitempty"`
	FinishedAt        *time.Time    `bun:"finished_at" json:"finished_at,omitempty"`
	Duration          time.Duration `bun:"duration" json:"duration"`
	DurationFormatted string        `bun:"duration_formatted" json:"duration_formatted"`
	TotalPaused       int8          `bun:"total_paused" json:"total_paused"`
}

func NewProcess(itemID uuid.UUID, processRuleID uuid.UUID) *Process {
	processCommonAttributes := ProcessCommonAttributes{
		ItemID:        itemID,
		ProcessRuleID: processRuleID,
		Status:        ProcessStatusPending,
	}

	return &Process{
		Entity:                  entity.NewEntity(),
		ProcessCommonAttributes: processCommonAttributes,
		ProcessTimeLogs:         ProcessTimeLogs{},
	}
}

func (p *Process) StartProcess(employeeID uuid.UUID) error {
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

func (p *Process) FinishProcess() error {
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

func (p *Process) PauseProcess() error {
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

func (p *Process) ContinueProcess() error {
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
