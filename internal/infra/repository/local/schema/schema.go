package schemarepositorylocal

import (
   "context"
   "sync"

   "github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type SchemaRepositoryLocal struct {
   mu      sync.RWMutex
   schemas map[string]model.Schema
}

func NewSchemaRepositoryLocal() model.SchemaRepository {
   return &SchemaRepositoryLocal{schemas: make(map[string]model.Schema)}
}

func (r *SchemaRepositoryLocal) NewSchema(ctx context.Context) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   // Initialize default schemas
   r.schemas[string(model.PUBLIC_SCHEMA)] = model.PUBLIC_SCHEMA
   r.schemas[string(model.LOST_SCHEMA)] = model.LOST_SCHEMA
   return nil
}

func (r *SchemaRepositoryLocal) UpdateSchema(ctx context.Context, p *model.Schema) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   r.schemas[string(*p)] = *p
   return nil
}

func (r *SchemaRepositoryLocal) DeleteSchema(ctx context.Context, id string) error {
   r.mu.Lock()
   defer r.mu.Unlock()
   delete(r.schemas, id)
   return nil
}

func (r *SchemaRepositoryLocal) GetSchemaById(ctx context.Context, id string) (*model.Schema, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   if s, ok := r.schemas[id]; ok {
       return &s, nil
   }
   return nil, nil
}

func (r *SchemaRepositoryLocal) GetSchemaByName(ctx context.Context, name string) (*model.Schema, error) {
   return r.GetSchemaById(ctx, name)
}

func (r *SchemaRepositoryLocal) GetAllSchemas(ctx context.Context) ([]model.Schema, error) {
   r.mu.RLock()
   defer r.mu.RUnlock()
   list := make([]model.Schema, 0, len(r.schemas))
   for _, s := range r.schemas {
       list = append(list, s)
   }
   return list, nil
}
