package tableentity

import "context"

type TableRepository interface {
	RegisterTable(ctx context.Context, Table *Table) error
	UpdateTable(ctx context.Context, Table *Table) error
	DeleteTable(ctx context.Context, id string) error
	GetTableById(ctx context.Context, id string) (*Table, error)
	GetAllTables(ctx context.Context) ([]Table, error)
}
