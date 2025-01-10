package orderqueuedto

import (
	"time"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type QueueDTO struct {
	GroupItemID       uuid.UUID     `json:"group_item_id"`
	ProcessRuleID     uuid.UUID     `json:"process_rule_id"`
	JoinedAt          time.Time     `json:"joined_at"`
	LeftAt            *time.Time    `json:"left_at"`
	Duration          time.Duration `json:"duration"`
	DurationFormatted string        `json:"duration_formatted"`
}

func (s *QueueDTO) FromDomain(queue *orderprocessentity.OrderQueue) {
	if queue == nil {
		return
	}
	*s = QueueDTO{
		GroupItemID:       queue.GroupItemID,
		ProcessRuleID:     queue.ProcessRuleID,
		JoinedAt:          queue.JoinedAt,
		LeftAt:            queue.LeftAt,
		Duration:          queue.Duration,
		DurationFormatted: queue.DurationFormatted,
	}
}
