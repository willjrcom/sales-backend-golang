package schemarepositorybun

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"golang.org/x/net/context"
)

type SchemaRepositoryBun struct {
	db *bun.DB
}

func NewSchemaRepositoryBun(db *bun.DB) *SchemaRepositoryBun {
	return &SchemaRepositoryBun{db: db}
}

func (r *SchemaRepositoryBun) NewSchema(ctx context.Context, id string) {
	database.LoadModels(ctx, r.db, id)
}
