package productcategorydto

import "github.com/google/uuid"

type CategoryMapDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ImagePath string    `json"image_path"`
}
