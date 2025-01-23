package clientrepositorylocal

import (
	"context"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ClientRepositoryLocal struct {
}

func NewClientRepositoryLocal() model.ClientRepository {
	return &ClientRepositoryLocal{}
}

func (r *ClientRepositoryLocal) CreateClient(ctx context.Context, p *model.Client) error {
	return nil
}

func (r *ClientRepositoryLocal) UpdateClient(ctx context.Context, p *model.Client) error {
	return nil
}

func (r *ClientRepositoryLocal) DeleteClient(ctx context.Context, id string) error {
	return nil
}

func (r *ClientRepositoryLocal) GetClientById(ctx context.Context, id string) (*model.Client, error) {
	return nil, nil
}

func (r *ClientRepositoryLocal) GetAllClients(ctx context.Context) ([]model.Client, error) {
	return nil, nil
}
