package groupitemdto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type GroupItemSnapshotDTO struct {
	ID          uuid.UUID       `json:"id,omitempty"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	GroupItemID uuid.UUID       `json:"group_item_id,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
}

func (s *GroupItemSnapshotDTO) FromDomain(snapshot *orderentity.GroupItemSnapshot) {
	if snapshot == nil {
		return
	}

	*s = GroupItemSnapshotDTO{
		ID:          snapshot.ID,
		CreatedAt:   snapshot.CreatedAt,
		GroupItemID: snapshot.GroupItemID,
		Data:        snapshot.Data,
	}
}
