package cliententity

import (
	"context"
)

type Repository interface {
	CreateClient(ctx context.Context, p *Client) error
	UpdateClient(ctx context.Context, p *Client) error
	DeleteClient(ctx context.Context, id string) error
	GetClientById(ctx context.Context, id string) (*Client, error)
	GetAllClients(ctx context.Context) ([]Client, error)
}
