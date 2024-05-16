package placedto

import (
	"errors"

	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

var (
	ErrNameRequired = errors.New("place name is required")
)

type CreatePlaceInput struct {
	tableentity.PlaceCommonAttributes
}

func (o *CreatePlaceInput) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *CreatePlaceInput) ToModel() (*tableentity.Place, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return tableentity.NewPlace(o.PlaceCommonAttributes), nil
}
