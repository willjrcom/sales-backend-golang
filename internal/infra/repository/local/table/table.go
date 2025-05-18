package tablerepositorylocal

import (
	"context"
	"sync"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type TableRepositoryLocal struct {
	mu     sync.RWMutex
	tables map[string]*model.Table
}

func NewTableRepositoryLocal() model.TableRepository {
	return &TableRepositoryLocal{tables: make(map[string]*model.Table)}
}

func (r *TableRepositoryLocal) CreateTable(ctx context.Context, table *model.Table) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tables == nil {
		r.tables = make(map[string]*model.Table)
	}
	r.tables[table.ID.String()] = table
	return nil
}

func (r *TableRepositoryLocal) UpdateTable(ctx context.Context, table *model.Table) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tables == nil {
		r.tables = make(map[string]*model.Table)
	}
	r.tables[table.ID.String()] = table
	return nil
}

func (r *TableRepositoryLocal) DeleteTable(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tables, id)
	return nil
}

func (r *TableRepositoryLocal) GetTableById(ctx context.Context, id string) (*model.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.tables == nil {
		return nil, nil
	}
	if tbl, ok := r.tables[id]; ok {
		return tbl, nil
	}
	return nil, nil
}

func (r *TableRepositoryLocal) GetAllTables(ctx context.Context) ([]model.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Table, 0, len(r.tables))
	for _, tbl := range r.tables {
		list = append(list, *tbl)
	}
	return list, nil
}

func (r *TableRepositoryLocal) GetUnusedTables(ctx context.Context) ([]model.Table, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Table, 0)
	for _, tbl := range r.tables {
		if tbl.IsAvailable {
			list = append(list, *tbl)
		}
	}
	return list, nil
}
