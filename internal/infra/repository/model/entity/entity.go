package entitymodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Entity struct {
	ID        uuid.UUID  `bun:"id,type:uuid,pk,notnull" json:"id"`
	CreatedAt time.Time  `bun:"created_at,notnull" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bun:"deleted_at" json:"deleted_at,omitempty"`
}

func NewEntity() Entity {
	now := time.Now().UTC()
	return Entity{ID: uuid.New(), CreatedAt: now, UpdatedAt: now}
}

func FromDomain(domain entity.Entity) Entity {
	return Entity{
		ID:        domain.ID,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
		DeletedAt: domain.DeletedAt,
	}
}

func (e *Entity) ToDomain() entity.Entity {
	return entity.Entity{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}
