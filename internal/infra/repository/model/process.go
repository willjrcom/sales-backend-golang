package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

var (
	ErrMustBeStarted = errors.New("process must be started")
	ErrMustBeReason  = errors.New("reason is required")
)

type OrderProcess struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_processes,alias:process"`
	OrderProcessTimeLogs
	OrderProcessCommonAttributes
}

type OrderProcessCommonAttributes struct {
	EmployeeID    *uuid.UUID    `bun:"employee_id,type:uuid"`
	GroupItemID   uuid.UUID     `bun:"group_item_id,type:uuid,notnull"`
	GroupItem     *GroupItem    `bun:"rel:belongs-to,join:group_item_id=id"`
	ProcessRuleID uuid.UUID     `bun:"process_rule_id,type:uuid,notnull"`
	Status        StatusProcess `bun:"status,notnull"`
	Products      []Product     `bun:"m2m:process_to_product_to_group_item,join:Process=Product"`
	Queue         *OrderQueue   `bun:"rel:has-one,join:group_item_id=group_item_id,process_rule_id=process_rule_id"`
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
