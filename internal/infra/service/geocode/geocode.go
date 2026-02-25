package geocodeservice

import (
	"context"
	"fmt"
	"os"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"googlemaps.github.io/maps"
)

func GetCoordinates(address *addressentity.AddressCommonAttributes) (*addressentity.Coordinates, error) {
	apiKey := os.Getenv("GCP_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GCP_KEY env var is not configured")
	}

	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Maps client: %v", err)
	}

	q := address.Street + " " + address.Number + ", " + address.Neighborhood + ", " + address.City + ", " + address.UF

	if address.Cep != "" {
		q += ", " + address.Cep + ", Brasil"
	} else {
		q += ", Brasil"
	}

	r := &maps.GeocodingRequest{
		Address: q,
	}

	results, err := c.Geocode(context.Background(), r)
	if err != nil {
		return nil, fmt.Errorf("geocoding error: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address")
	}

	location := results[0].Geometry.Location
	return &addressentity.Coordinates{Latitude: location.Lat, Longitude: location.Lng}, nil
}
