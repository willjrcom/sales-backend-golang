package schemarepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"golang.org/x/net/context"
)

type SchemaRepositoryBun struct {
	db *bun.DB
	mu sync.Mutex
}

func NewSchemaRepositoryBun(db *bun.DB) *SchemaRepositoryBun {
	return &SchemaRepositoryBun{db: db}
}

func (r *SchemaRepositoryBun) NewSchema(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.RegisterModels(ctx, r.db); err != nil {
		return err
	}

	if err := database.LoadCompanyModels(ctx, r.db); err != nil {
		return err
	}

	if err := setupFtSearch(ctx, r.db); err != nil {
		return err
	}

	return nil
}
