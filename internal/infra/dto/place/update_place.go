package placedto

import (
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type UpdatePlaceInput struct {
	tableentity.PatchPlace
}

func (c *UpdatePlaceInput) UpdateModel(place *tableentity.Place) (err error) {
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
