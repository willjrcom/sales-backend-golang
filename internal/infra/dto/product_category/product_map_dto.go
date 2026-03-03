package productcategorydto

import "github.com/google/uuid"

type ProductMapDTO struct {
	VariationID uuid.UUID `json:"variation_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
}
