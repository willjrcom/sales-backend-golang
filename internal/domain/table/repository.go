package tableentity

import "context"

type TableRepository interface {
	CreateTable(ctx context.Context, table *Table) error
	UpdateTable(ctx context.Context, table *Table) error
	DeleteTable(ctx context.Context, id string) error
	GetTableById(ctx context.Context, id string) (*Table, error)
	GetAllTables(ctx context.Context) ([]Table, error)
}

type PlaceRepository interface {
	CreatePlace(ctx context.Context, place *Place) error
	UpdatePlace(ctx context.Context, place *Place) error
	DeletePlace(ctx context.Context, id string) error
	GetPlaceById(ctx context.Context, id string) (*Place, error)
	GetAllPlaces(ctx context.Context) ([]Place, error)
}
