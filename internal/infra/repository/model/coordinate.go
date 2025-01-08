package model

import addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func (c *Coordinates) FromDomain(coordinates *addressentity.Coordinates) {
	if coordinates == nil {
		return
	}
	*c = Coordinates{
		Latitude:  coordinates.Latitude,
		Longitude: coordinates.Longitude,
	}
}

func (c *Coordinates) ToDomain() *addressentity.Coordinates {
	return &addressentity.Coordinates{
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
	}
}
