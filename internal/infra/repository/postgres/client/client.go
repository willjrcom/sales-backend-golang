package clientrepositorylocal

import (
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	"golang.org/x/net/context"
)

type ProductRepositoryBun struct {
}

func NewProductRepositoryLocal() *ProductRepositoryBun {
	return &ProductRepositoryBun{}
}

func (r *ProductRepositoryBun) RegisterClient(ctx context.Context, c *cliententity.Client) error {
	return nil
}

func (r *ProductRepositoryBun) UpdateClient(ctx context.Context, c *cliententity.Client) error {
	return nil
}

func (r *ProductRepositoryBun) DeleteClient(ctx context.Context, id string) error {
	return nil
}

func (r *ProductRepositoryBun) GetClientById(ctx context.Context, id string) (*cliententity.Client, error) {
	return nil, nil
}

func (r *ProductRepositoryBun) GetClientBy(ctx context.Context, c *cliententity.Client) (*cliententity.Client, error) {
	return nil, nil
}

func (r *ProductRepositoryBun) GetAllClient(ctx context.Context) ([]cliententity.Client, error) {
	return nil, nil
}
