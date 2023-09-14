package clientrepositorylocal

import (
	"context"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
)

type ClientRepositoryLocal struct {
}

func NewClientRepositoryLocal() *ClientRepositoryLocal {
	return &ClientRepositoryLocal{}
}

func (r *ClientRepositoryLocal) RegisterClient(ctx context.Context, p *cliententity.Client) error {
	return nil
}

func (r *ClientRepositoryLocal) UpdateClient(ctx context.Context, p *cliententity.Client) error {
	return nil
}

func (r *ClientRepositoryLocal) DeleteClient(ctx context.Context, id string) error {
	return nil
}

func (r *ClientRepositoryLocal) GetClientById(ctx context.Context, id string) (*cliententity.Client, error) {
	return nil, nil
}

func (r *ClientRepositoryLocal) GetClientBy(ctx context.Context, c *cliententity.Client) ([]cliententity.Client, error) {
	return nil, nil
}

func (r *ClientRepositoryLocal) GetAllClient(ctx context.Context) ([]cliententity.Client, error) {
	return nil, nil
}
