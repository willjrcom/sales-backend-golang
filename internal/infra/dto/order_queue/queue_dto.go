package orderqueuedto

import (
	"time"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type QueueDTO struct {
	GroupItemID       uuid.UUID
	ProcessRuleID     *uuid.UUID
	JoinedAt          time.Time
	LeftAt            *time.Time
	Duration          time.Duration
	DurationFormatted string
}

func (s *QueueDTO) FromDomain(queue *orderprocessentity.OrderQueue) {
	*s = QueueDTO{
		GroupItemID:       queue.GroupItemID,
		ProcessRuleID:     queue.ProcessRuleID,
		JoinedAt:          queue.JoinedAt,
		LeftAt:            queue.LeftAt,
		Duration:          queue.Duration,
		DurationFormatted: queue.DurationFormatted,
	}
}
