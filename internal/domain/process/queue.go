package processentity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Queue struct {
	entity.Entity
	bun.BaseModel `bun:"table:queues"`
	QueueCommonAttributes
	QueueTimeLogs
}

type QueueCommonAttributes struct {
	ItemID    uuid.UUID `bun:"column:item_id,type:uuid,notnull" json:"item_id"`
	ProductID uuid.UUID `bun:"column:product_id,type:uuid,notnull" json:"product_id"`
}

type QueueTimeLogs struct {
	JoinedAt          time.Time     `bun:"joined_at" json:"joined_at,omitempty"`
	LeftAt            *time.Time    `bun:"left_at" json:"left_at,omitempty"`
	Duration          time.Duration `bun:"duration" json:"duration"`
	DurationFormatted string        `bun:"duration_formatted" json:"duration_formatted"`
}

func NewQueue(process *Process) (*Queue, error) {
	if process.FinishedAt == nil {
		return nil, errors.New("process must be finished")
	}

	return &Queue{
		Entity: entity.NewEntity(),
		QueueCommonAttributes: QueueCommonAttributes{
			ItemID:    process.GroupItemID,
			ProductID: process.ProductID,
		},
		QueueTimeLogs: QueueTimeLogs{
			JoinedAt: *process.FinishedAt,
		},
	}, nil
}

func (q *Queue) LeftQueue(nextProcess *Process) error {
	if !q.JoinedAt.Before(*nextProcess.StartedAt) {
		return errors.New("next process must be started after queue joined")
	}

	q.LeftAt = &time.Time{}
	*q.LeftAt = *nextProcess.StartedAt
	q.Duration = nextProcess.StartedAt.Sub(q.JoinedAt)
	q.DurationFormatted = q.Duration.String()
	return nil
}
