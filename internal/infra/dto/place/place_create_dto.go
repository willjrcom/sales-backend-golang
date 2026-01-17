package placedto

import (
	"errors"

	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

var (
	ErrNameRequired = errors.New("place name is required")
)

type CreatePlaceInput struct {
	Name        string  `json:"name"`
	ImagePath   *string `json:"image_path"`
	IsAvailable bool    `json:"is_available"`
	IsActive    *bool   `json:"is_active"`
}

func (o *CreatePlaceInput) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *CreatePlaceInput) ToDomain() (*tableentity.Place, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	isActive := true
	if o.IsActive != nil {
		isActive = *o.IsActive
	}
	placeCommonAttributes := tableentity.PlaceCommonAttributes{
		Name:        o.Name,
		ImagePath:   o.ImagePath,
		IsAvailable: o.IsAvailable,
		IsActive:    isActive,
	}

	return tableentity.NewPlace(placeCommonAttributes), nil
}
