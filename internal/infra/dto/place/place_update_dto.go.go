package placedto

import (
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type PlaceUpdateDTO struct {
	Name        *string `json:"name"`
	ImagePath   *string `json:"image_path"`
	IsAvailable *bool   `json:"is_available"`
}

func (c *PlaceUpdateDTO) UpdateDomain(place *tableentity.Place) (err error) {
	if c.Name != nil {
		place.Name = *c.Name
	}

	if c.ImagePath != nil {
		place.ImagePath = c.ImagePath
	}

	if c.IsAvailable != nil {
		place.IsAvailable = *c.IsAvailable
	}

	return nil
}
