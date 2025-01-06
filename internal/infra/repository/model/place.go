package model

import (
	"github.com/uptrace/bun"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Place struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:places"`
	PlaceCommonAttributes
}

type PlaceCommonAttributes struct {
	Name        string          `bun:"name,notnull"`
	ImagePath   *string         `bun:"image_path"`
	IsAvailable bool            `bun:"is_available"`
	Tables      []PlaceToTables `bun:"rel:has-many,join:id=place_id"`
}

func (p *Place) FromDomain(place *tableentity.Place) {
	*p = Place{
		Entity: entitymodel.FromDomain(place.Entity),
		PlaceCommonAttributes: PlaceCommonAttributes{
			Name:        place.Name,
			ImagePath:   place.ImagePath,
			IsAvailable: place.IsAvailable,
			Tables:      []PlaceToTables{},
		},
	}

	for _, table := range place.Tables {
		pt := PlaceToTables{}
		pt.FromDomain(&table)
		p.Tables = append(p.Tables, pt)
	}
}

func (p *Place) ToDomain() *tableentity.Place {
	if p == nil {
		return nil
	}
	place := &tableentity.Place{
		Entity: p.Entity.ToDomain(),
		PlaceCommonAttributes: tableentity.PlaceCommonAttributes{
			Name:        p.Name,
			ImagePath:   p.ImagePath,
			IsAvailable: p.IsAvailable,
			Tables:      []tableentity.PlaceToTables{},
		},
	}

	for _, table := range p.Tables {
		place.Tables = append(place.Tables, *table.ToDomain())
	}

	return place
}
