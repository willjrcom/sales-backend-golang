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
