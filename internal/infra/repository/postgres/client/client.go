package clientrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"

	"golang.org/x/net/context"
)

type ClientRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewClientRepositoryBun(db *bun.DB) *ClientRepositoryBun {
	return &ClientRepositoryBun{db: db}
}

func (r *ClientRepositoryBun) RegisterClient(ctx context.Context, c *cliententity.Client) error {
	r.mu.Lock()
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Register client
	if _, err := tx.NewInsert().Model(c).Exec(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// Register contact
	if _, err := tx.NewInsert().Model(&c.Contact).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	// Register addresse
	if _, err := tx.NewInsert().Model(&c.Address).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	if err := tx.Commit(); err != nil {
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
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Delete client
	if _, err = tx.NewDelete().Model(&cliententity.Client{}).Where("id = ?", id).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	// Delete contact
	if _, err = tx.NewDelete().Model(&personentity.Contact{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	// Delete addresse
	if _, err = tx.NewDelete().Model(&addressentity.Address{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		return rollback(&tx, err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) GetClientById(ctx context.Context, id string) (*cliententity.Client, error) {
	client := &cliententity.Client{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(client).Where("client.id = ?", id).Relation("Address").Relation("Contact").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *ClientRepositoryBun) GetAllClients(ctx context.Context) ([]cliententity.Client, error) {
	clients := []cliententity.Client{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&clients).Relation("Address").Relation("Contact").Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return clients, nil
}

func rollback(tx *bun.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		return err
	}

	return err
}
