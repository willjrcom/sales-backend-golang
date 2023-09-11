package entity

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID        uuid.UUID  `bun:"id,type:uuid,pk"`
	CreatedAt time.Time  `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at"`
}

func NewEntity() Entity {
	return Entity{ID: uuid.New(), CreatedAt: time.Now()}
}
