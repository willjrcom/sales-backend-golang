package entity

import (
	"time"

	"github.com/google/uuid"
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
