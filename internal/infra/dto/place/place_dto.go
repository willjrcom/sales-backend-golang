package placedto

import (
	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type PlaceDTO struct {
	ID uuid.UUID `json:"id"`
	PlaceCommonAttributes
}

type PlaceCommonAttributes struct {
	Name        string             `json:"name"`
	ImagePath   *string            `json:"image_path"`
	IsAvailable bool               `json:"is_available"`
	IsActive    bool               `json:"is_active"`
	Tables      []PlaceToTablesDTO `json:"tables"`
}

func (p *PlaceDTO) FromDomain(place *tableentity.Place) {
	if place == nil {
		return
	}
	*p = PlaceDTO{
		ID: place.ID,
		PlaceCommonAttributes: PlaceCommonAttributes{
			Name:        place.Name,
			ImagePath:   place.ImagePath,
			IsAvailable: place.IsAvailable,
			IsActive:    place.IsActive,
			Tables:      []PlaceToTablesDTO{},
		},
	}

	for _, table := range place.Tables {
		pt := PlaceToTablesDTO{}
		pt.FromDomain(&table)
		p.Tables = append(p.Tables, pt)
	}

	if len(place.Tables) == 0 {
		p.Tables = nil
	}
}
