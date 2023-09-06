package productentity

import (
	"github.com/google/uuid"
)

type CategoryProduct struct {
	ID       uuid.UUID
	Name     string
	Sizes    []string
	Products []*Product `bun:"rel:has-many"`
}
