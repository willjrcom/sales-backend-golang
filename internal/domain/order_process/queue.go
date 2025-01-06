package orderprocessentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type OrderQueue struct {
	entity.Entity
	OrderQueueCommonAttributes
	OrderQueueTimeLogs
}

type OrderQueueCommonAttributes struct {
	GroupItemID   uuid.UUID
	ProcessRuleID uuid.UUID
}

type OrderQueueTimeLogs struct {
	JoinedAt          time.Time
	LeftAt            *time.Time
	Duration          time.Duration
	DurationFormatted string
}

func NewOrderQueue(groupItemID uuid.UUID, joinedAt time.Time) (*OrderQueue, error) {
	return &OrderQueue{
		Entity: entity.NewEntity(),
		OrderQueueCommonAttributes: OrderQueueCommonAttributes{
			GroupItemID:   groupItemID,
			ProcessRuleID: uuid.Nil,
		},
		OrderQueueTimeLogs: OrderQueueTimeLogs{
			JoinedAt:          joinedAt,
			DurationFormatted: "0s",
		},
	}, nil
}

func (q *OrderQueue) Finish(processRuleID uuid.UUID, finishedAt time.Time) {
	q.ProcessRuleID = processRuleID
	q.LeftAt = &time.Time{}
	*q.LeftAt = finishedAt
	q.Duration = q.LeftAt.Sub(q.JoinedAt)
	q.DurationFormatted = q.Duration.String()
}
