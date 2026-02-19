package placedto

import orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"

type PlaceUpdateDTO struct {
	Name        *string `json:"name"`
	ImagePath   *string `json:"image_path"`
	IsAvailable *bool   `json:"is_available"`
	IsActive    *bool   `json:"is_active"`
}

func (c *PlaceUpdateDTO) UpdateDomain(place *orderentity.Place) (err error) {
	if c.Name != nil {
		place.Name = *c.Name
	}

	if c.ImagePath != nil {
		place.ImagePath = c.ImagePath
	}

	if c.IsAvailable != nil {
		place.IsAvailable = *c.IsAvailable
	}

	if c.IsActive != nil {
		place.IsActive = *c.IsActive
	}

	return nil
}
