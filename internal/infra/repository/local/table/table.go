package tablerepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type TableRepositoryLocal struct {}

func NewTableRepositoryLocal() model.TableRepository {
	return &TableRepositoryLocal{}
}

func (r *TableRepositoryLocal) CreateTable(ctx context.Context, table *model.Table) error {
	return nil
}

func (r *TableRepositoryLocal) UpdateTable(ctx context.Context, table *model.Table) error {
	return nil
}

func (r *TableRepositoryLocal) DeleteTable(ctx context.Context, id string) error {
	return nil
}

func (r *TableRepositoryLocal) GetTableById(ctx context.Context, id string) (*model.Table, error) {
	return nil, nil
}

func (r *TableRepositoryLocal) GetAllTables(ctx context.Context) ([]model.Table, error) {
	return nil, nil
}
