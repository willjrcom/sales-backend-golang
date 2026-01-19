package productcategorydto

import "github.com/google/uuid"

type CategoryMapDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
