package clientrepositorylocal

import (
	"sync"

	"github.com/uptrace/bun"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	"golang.org/x/net/context"
)

type ClientRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewClientRepositoryLocal(db *bun.DB) *ClientRepositoryBun {
	return &ClientRepositoryBun{db: db}
}

func (r *ClientRepositoryBun) RegisterClient(ctx context.Context, c *cliententity.Client) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(c).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) UpdateClient(ctx context.Context, c *cliententity.Client) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) DeleteClient(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&cliententity.Client{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) GetClientById(ctx context.Context, id string) (*cliententity.Client, error) {
	client := &cliententity.Client{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(client).Where("client.id = ?", id).Relation("Address").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *ClientRepositoryBun) GetClientBy(ctx context.Context, c *cliententity.Client) ([]cliententity.Client, error) {
	clients := []cliententity.Client{}

	r.mu.Lock()
	query := r.db.NewSelect().Model(&cliententity.Client{})

	if c.Name != "" {
		query.Where("client.code = ?", c.Name)
	}

	err := query.Relation("Address").Scan(ctx, &clients)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepositoryBun) GetAllClient(ctx context.Context) ([]cliententity.Client, error) {
	clients := []cliententity.Client{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&clients).Relation("Address").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return clients, nil
}
