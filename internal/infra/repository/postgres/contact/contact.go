package contactrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"golang.org/x/net/context"
)

type ContactRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewContactRepositoryBun(db *bun.DB) model.ContactRepository {
	return &ContactRepositoryBun{db: db}
}

func (r *ContactRepositoryBun) CreateContact(ctx context.Context, c *model.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	_, err := r.db.NewInsert().Model(c).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepositoryBun) UpdateContact(ctx context.Context, c *model.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	_, err := r.db.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepositoryBun) DeleteContact(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	_, err := r.db.NewDelete().Model(&model.Contact{}).Where("id = ?", id).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepositoryBun) GetContactById(ctx context.Context, id string) (*model.Contact, error) {
	contact := &model.Contact{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	err := r.db.NewSelect().Model(contact).Where("contact.id = ?", id).Scan(ctx)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *ContactRepositoryBun) GetContactByDddAndNumber(ctx context.Context, ddd string, number string, contactType string) (*model.Contact, error) {
	contact := &model.Contact{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(contact).Where("ddd = ? AND number = ? AND type = ?", ddd, number, contactType).Scan(ctx); err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *ContactRepositoryBun) FtSearchContacts(ctx context.Context, text string, contactType string) (contacts []model.Contact, err error) {
	contacts = []model.Contact{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	err = r.db.NewSelect().Model(&contacts).Where("ts @@ websearch_to_tsquery('simple', ?) and type = ?", text, contactType).Scan(ctx)
	return
}
