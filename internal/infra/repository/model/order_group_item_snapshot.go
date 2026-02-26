package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type OrderGroupItemSnapshot struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:order_group_item_snapshots"`
	GroupItemID   uuid.UUID       `bun:"column:group_item_id,type:uuid,notnull,unique"`
	Data          json.RawMessage `bun:"column:data,type:jsonb,notnull"`
}

func (s *OrderGroupItemSnapshot) FromDomain(snapshot *orderentity.GroupItemSnapshot) {
	if snapshot == nil {
		return
	}
	s.Entity = entitymodel.FromDomain(snapshot.Entity)
	s.GroupItemID = snapshot.GroupItemID
	s.Data = snapshot.Data
}

func (s *OrderGroupItemSnapshot) ToDomain() *orderentity.GroupItemSnapshot {
	if s == nil {
		return nil
	}
	return &orderentity.GroupItemSnapshot{
		Entity:      s.Entity.ToDomain(),
		GroupItemID: s.GroupItemID,
		Data:        s.Data,
	}
}
