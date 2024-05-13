package orderprocessentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type OrderQueue struct {
	entity.Entity
	bun.BaseModel `bun:"table:order_queues"`
	OrderQueueCommonAttributes
	OrderQueueTimeLogs
}

type OrderQueueCommonAttributes struct {
	GroupItemID   uuid.UUID  `bun:"column:group_item_id,type:uuid,notnull" json:"group_item_id"`
	ProcessRuleID *uuid.UUID `bun:"column:process_rule_id,type:uuid,notnull" json:"process_rule_id,omitempty"`
}

type OrderQueueTimeLogs struct {
	JoinedAt          time.Time     `bun:"joined_at" json:"joined_at,omitempty"`
	LeftAt            *time.Time    `bun:"left_at" json:"left_at,omitempty"`
	Duration          time.Duration `bun:"duration" json:"duration"`
	DurationFormatted string        `bun:"duration_formatted" json:"duration_formatted"`
}

func NewOrderQueue(groupItemID uuid.UUID, joinedAt time.Time) (*OrderQueue, error) {
	return &OrderQueue{
		Entity: entity.NewEntity(),
		OrderQueueCommonAttributes: OrderQueueCommonAttributes{
			GroupItemID: groupItemID,
		},
		OrderQueueTimeLogs: OrderQueueTimeLogs{
			JoinedAt:          joinedAt,
			DurationFormatted: "0s",
		},
	}, nil
}

func (q *OrderQueue) Finish(processRuleID uuid.UUID, finishedAt time.Time) {
	q.ProcessRuleID = &processRuleID
	q.LeftAt = &time.Time{}
	*q.LeftAt = finishedAt
	q.Duration = q.LeftAt.Sub(q.JoinedAt)
	q.DurationFormatted = q.Duration.String()
}
