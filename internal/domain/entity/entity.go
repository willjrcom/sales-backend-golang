package entity

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewEntity() Entity {
	return Entity{ID: uuid.New(), CreatedAt: time.Now()}
}
