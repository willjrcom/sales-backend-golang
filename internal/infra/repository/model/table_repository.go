package model

import (
	"context"

	"github.com/google/uuid"
)

type TableRepository interface {
	CreateTable(ctx context.Context, table *Table) error
	UpdateTable(ctx context.Context, table *Table) error
	DeleteTable(ctx context.Context, id string) error
	GetTableById(ctx context.Context, id string) (*Table, error)
	GetAllTables(ctx context.Context) ([]Table, error)
	GetUnusedTables(ctx context.Context) ([]Table, error)
}

type PlaceRepository interface {
	CreatePlace(ctx context.Context, place *Place) error
	UpdatePlace(ctx context.Context, place *Place) error
	DeletePlace(ctx context.Context, id string) error
	GetPlaceById(ctx context.Context, id string) (*Place, error)
	GetAllPlaces(ctx context.Context) ([]Place, error)
	AddTableToPlace(ctx context.Context, placeToTables *PlaceToTables) error
	GetTableToPlaceByPlaceIDAndPosition(ctx context.Context, placeID uuid.UUID, column, row int) (*PlaceToTables, error)
	GetTableToPlaceByTableID(ctx context.Context, table uuid.UUID) (*PlaceToTables, error)
	RemoveTableFromPlace(ctx context.Context, tableID uuid.UUID) error
}
