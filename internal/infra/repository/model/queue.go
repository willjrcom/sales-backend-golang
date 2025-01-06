package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type OrderQueue struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_queues"`
	OrderQueueCommonAttributes
	OrderQueueTimeLogs
}

type OrderQueueCommonAttributes struct {
	GroupItemID   uuid.UUID  `bun:"column:group_item_id,type:uuid,notnull"`
	ProcessRuleID *uuid.UUID `bun:"column:process_rule_id,type:uuid"`
}

type OrderQueueTimeLogs struct {
	JoinedAt          time.Time     `bun:"joined_at"`
	LeftAt            *time.Time    `bun:"left_at"`
	Duration          time.Duration `bun:"duration"`
	DurationFormatted string        `bun:"duration_formatted"`
}
