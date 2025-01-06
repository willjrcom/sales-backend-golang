package placedto

import (
	"errors"

	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

var (
	ErrTableIDRequired = errors.New("table id is required")
	ErrPlaceIDRequired = errors.New("place id is required")
	ErrInvalidColumn   = errors.New("column must be positive")
	ErrInvalidRow      = errors.New("row must be positive")
)

type PlaceUpdateTableDTO struct {
	TableID uuid.UUID `json:"table_id"`
	PlaceID uuid.UUID `json:"place_id"`
	Column  int       `json:"column"`
	Row     int       `json:"row"`
}

func (o *PlaceUpdateTableDTO) validate() error {
	if o.TableID == uuid.Nil {
		return ErrTableIDRequired
	}

	if o.PlaceID == uuid.Nil {
		return ErrPlaceIDRequired
	}

	if o.Column < 0 {
		return ErrInvalidColumn
	}

	if o.Row < 0 {
		return ErrInvalidRow
	}
	return nil
}

func (o *PlaceUpdateTableDTO) ToDomain() (placeToTables *tableentity.PlaceToTables, err error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	return tableentity.NewPlaceToTable(o.PlaceID, o.TableID, o.Column, o.Row), nil
}
