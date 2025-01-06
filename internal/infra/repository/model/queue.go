package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
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

func (q *OrderQueue) FromDomain(queue *orderprocessentity.OrderQueue) {
	*q = OrderQueue{
		Entity: entitymodel.FromDomain(queue.Entity),
		OrderQueueCommonAttributes: OrderQueueCommonAttributes{
			GroupItemID:   queue.GroupItemID,
			ProcessRuleID: queue.ProcessRuleID,
		},
		OrderQueueTimeLogs: OrderQueueTimeLogs{
			JoinedAt:          queue.JoinedAt,
			LeftAt:            queue.LeftAt,
			Duration:          queue.Duration,
			DurationFormatted: queue.DurationFormatted,
		},
	}
}

func (q *OrderQueue) ToDomain() *orderprocessentity.OrderQueue {
	return &orderprocessentity.OrderQueue{
		Entity: q.Entity.ToDomain(),
		OrderQueueCommonAttributes: orderprocessentity.OrderQueueCommonAttributes{
			GroupItemID:   q.GroupItemID,
			ProcessRuleID: q.ProcessRuleID,
		},
		OrderQueueTimeLogs: orderprocessentity.OrderQueueTimeLogs{
			JoinedAt:          q.JoinedAt,
			LeftAt:            q.LeftAt,
			Duration:          q.Duration,
			DurationFormatted: q.DurationFormatted,
		},
	}
}
