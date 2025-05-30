package model

import (
	"context"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, p *Client) error
	UpdateClient(ctx context.Context, p *Client) error
	DeleteClient(ctx context.Context, id string) error
   GetClientById(ctx context.Context, id string) (*Client, error)
   // GetAllClients returns a paginated list of clients and the total count.
   GetAllClients(ctx context.Context, offset, limit int) ([]Client, int, error)
}
