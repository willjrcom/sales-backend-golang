package orderentity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type GroupItemSnapshot struct {
	entity.Entity
	GroupItemID uuid.UUID
	Data        []byte
}
