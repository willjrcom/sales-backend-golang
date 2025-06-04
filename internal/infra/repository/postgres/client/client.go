package clientrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"

	"golang.org/x/net/context"
)

type ClientRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewClientRepositoryBun(db *bun.DB) model.ClientRepository {
	return &ClientRepositoryBun{db: db}
}

func (r *ClientRepositoryBun) CreateClient(ctx context.Context, c *model.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Create client
	if _, err := tx.NewInsert().Model(c).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	// Create contact
	if _, err := tx.NewInsert().Model(c.Contact).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	// Create addresse
	if _, err := tx.NewInsert().Model(c.Address).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) UpdateClient(ctx context.Context, c *model.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		return err
	}

	if c.Contact != nil {
		if _, err := tx.NewDelete().Model(&model.Contact{}).Where("object_id = ?", c.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create contact
		if _, err := tx.NewInsert().Model(c.Contact).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if c.Address != nil {
		if _, err := tx.NewDelete().Model(&model.Address{}).Where("object_id = ?", c.ID).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		// Create addresse
		if _, err := tx.NewInsert().Model(c.Address).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) DeleteClient(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	// Delete client
	if _, err = tx.NewDelete().Model(&model.Client{}).Where("id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	// Delete contact
	if _, err = tx.NewDelete().Model(&model.Contact{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	// Delete addresse
	if _, err = tx.NewDelete().Model(&model.Address{}).Where("object_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) GetClientById(ctx context.Context, id string) (*model.Client, error) {
	client := &model.Client{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(client).Where("client.id = ?", id).Relation("Address").Relation("Contact").Scan(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

// GetAllClients retrieves a paginated list of clients and the total count.
func (r *ClientRepositoryBun) GetAllClients(ctx context.Context, page, perPage int) ([]model.Client, int, error) {
	var clients []model.Client
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, 0, err
	}

	// count total records
	totalCount, err := r.db.NewSelect().Model((*model.Client)(nil)).Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	// fetch paginated records
		if err := r.db.NewSelect().
			Model(&clients).
			Relation("Address").
			Relation("Contact").
			Limit(perPage).
			Offset(page * perPage).
			Scan(ctx); err != nil {
		return nil, 0, err
	}
	return clients, int(totalCount), nil
}
