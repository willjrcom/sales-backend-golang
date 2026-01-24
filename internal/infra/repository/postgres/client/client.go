package clientrepositorybun

import (
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"

	"golang.org/x/net/context"
)

type ClientRepositoryBun struct {
	db *bun.DB
}

func NewClientRepositoryBun(db *bun.DB) model.ClientRepository {
	return &ClientRepositoryBun{db: db}
}

func (r *ClientRepositoryBun) CreateClient(ctx context.Context, c *model.Client) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Create client
	if _, err := tx.NewInsert().Model(c).Exec(ctx); err != nil {
		return err
	}

	// Create contact
	if _, err := tx.NewInsert().Model(c.Contact).Exec(ctx); err != nil {
		return err
	}

	// Create addresse
	if _, err := tx.NewInsert().Model(c.Address).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) UpdateClient(ctx context.Context, c *model.Client) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		return err
	}

	if c.Contact != nil {
		res, err := tx.NewUpdate().Model(c.Contact).Where("id = ?", c.Contact.ID).Exec(ctx)
		if err != nil {
			return err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			if _, err := tx.NewInsert().Model(c.Contact).Exec(ctx); err != nil {
				return err
			}
		}
	}

	if c.Address != nil {
		res, err := tx.NewUpdate().Model(c.Address).Where("id = ?", c.Address.ID).Exec(ctx)
		if err != nil {
			return err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			if _, err := tx.NewInsert().Model(c.Address).Exec(ctx); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) DeleteClient(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Soft delete: set is_active to false on client
	isActive := false
	if _, err = tx.NewUpdate().
		Model(&model.Client{}).
		Set("is_active = ?", isActive).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepositoryBun) GetClientById(ctx context.Context, id string) (*model.Client, error) {
	client := &model.Client{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(client).Where("client.id = ?", id).Relation("Address").Relation("Contact").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return client, nil
}

// GetAllClients retrieves a paginated list of clients and the total count.
func (r *ClientRepositoryBun) GetAllClients(ctx context.Context, page, perPage int, isActive ...bool) ([]model.Client, int, error) {
	var clients []model.Client

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	// Default to active records (true)
	activeFilter := true
	if len(isActive) > 0 {
		activeFilter = isActive[0]
	}

	// count total records
	totalCount, err := tx.NewSelect().Model((*model.Client)(nil)).Where("is_active = ?", activeFilter).Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	// fetch paginated records
	if err := tx.NewSelect().
		Model(&clients).
		Where("client.is_active = ?", activeFilter).
		Relation("Address").
		Relation("Contact").
		Limit(perPage).
		Offset(page * perPage).
		Scan(ctx); err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return clients, int(totalCount), nil
}
