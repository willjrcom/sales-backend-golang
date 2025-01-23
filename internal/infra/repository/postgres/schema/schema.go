package schemarepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"golang.org/x/net/context"
)

type SchemaRepositoryBun struct {
	db *bun.DB
	mu sync.Mutex
}

func NewSchemaRepositoryBun(db *bun.DB) model.SchemaRepository {
	return &SchemaRepositoryBun{db: db}
}

func (r *SchemaRepositoryBun) NewSchema(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.CreateNewCompanySchema(ctx, r.db); err != nil {
		return err
	}

	return nil
}
