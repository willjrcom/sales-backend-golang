package orderprocessentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrMustBeStarted = errors.New("process must be started")
	ErrMustBeReason  = errors.New("reason is required")
)

type OrderProcess struct {
	entity.Entity
	OrderProcessTimeLogs
	OrderProcessCommonAttributes
}

type OrderProcessCommonAttributes struct {
	EmployeeID    *uuid.UUID
	GroupItemID   uuid.UUID
	GroupItem     *orderentity.GroupItem
	ProcessRuleID uuid.UUID
	Status        StatusProcess
	Products      []productentity.Product
	Queue         *OrderQueue
}

type OrderProcessTimeLogs struct {
	StartedAt         *time.Time
	PausedAt          *time.Time
	ContinuedAt       *time.Time
	FinishedAt        *time.Time
	CanceledAt        *time.Time
	CanceledReason    *string
	Duration          time.Duration
	DurationFormatted string
	TotalPaused       int8
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
	*p.StartedAt = time.Now().UTC()
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
	*p.FinishedAt = time.Now().UTC()
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
	*p.PausedAt = time.Now().UTC()
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
	*p.ContinuedAt = time.Now().UTC()
	p.Status = ProcessStatusContinued
	p.PausedAt = nil
	return nil
}

func (p *OrderProcess) CancelProcess(reason *string) error {
	if reason == nil {
		return ErrMustBeReason
	}

	p.CanceledReason = reason
	p.CanceledAt = &time.Time{}
	*p.CanceledAt = time.Now().UTC()
	p.Status = ProcessStatusCanceled
	return nil
}
