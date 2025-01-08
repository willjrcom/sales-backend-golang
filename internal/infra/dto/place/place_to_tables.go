package placedto

import (
	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
	tabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table"
)

type PlaceToTablesDTO struct {
	PlaceID uuid.UUID          `json:"place_id"`
	Place   *PlaceDTO          `json:"place,omitempty"`
	TableID uuid.UUID          `json:"table_id"`
	Table   *tabledto.TableDTO `json:"table"`
	Column  int                `json:"column"`
	Row     int                `json:"row"`
}

func (p *PlaceToTablesDTO) FromDomain(placeToTables *tableentity.PlaceToTables) {
	if placeToTables == nil {
		return
	}
	*p = PlaceToTablesDTO{
		PlaceID: placeToTables.PlaceID,
		Place:   &PlaceDTO{},
		TableID: placeToTables.TableID,
		Table:   &tabledto.TableDTO{},
		Column:  placeToTables.Column,
		Row:     placeToTables.Row,
	}

	p.Place.FromDomain(placeToTables.Place)
	p.Table.FromDomain(placeToTables.Table)
}
