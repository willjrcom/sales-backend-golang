package schemarepositorylocal

import (
	"context"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type SchemaRepositoryLocal struct{}

func NewSchemaRepositoryLocal() model.SchemaRepository {
	return &SchemaRepositoryLocal{}
}

func (r *SchemaRepositoryLocal) NewSchema(ctx context.Context) error {
	return nil
}

func (r *SchemaRepositoryLocal) UpdateSchema(ctx context.Context, p *model.Schema) error {
	return nil
}

func (r *SchemaRepositoryLocal) DeleteSchema(ctx context.Context, id string) error {
	return nil
}

func (r *SchemaRepositoryLocal) GetSchemaById(ctx context.Context, id string) (*model.Schema, error) {
	return nil, nil
}

func (r *SchemaRepositoryLocal) GetSchemaByName(ctx context.Context, name string) (*model.Schema, error) {
	return nil, nil
}

func (r *SchemaRepositoryLocal) GetAllSchemas(ctx context.Context) ([]model.Schema, error) {
	return nil, nil
}
