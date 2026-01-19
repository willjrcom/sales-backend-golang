package productcategorydto

import "github.com/google/uuid"

type ProductMapDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
