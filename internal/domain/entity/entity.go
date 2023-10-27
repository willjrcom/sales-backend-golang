package entity

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID        uuid.UUID  `bun:"id,type:uuid,pk,notnull"`
	CreatedAt time.Time  `bun:"created_at,notnull"`
	UpdatedAt time.Time  `bun:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at"`
}

func NewEntity() Entity {
	now := time.Now()
	return Entity{ID: uuid.New(), CreatedAt: now, UpdatedAt: now}
}
