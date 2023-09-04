package productEntity

import "github.com/google/uuid"

type CategoryProduct struct {
	ID    uuid.UUID
	Name  string
	Sizes []string
}
