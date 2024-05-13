package contactrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"golang.org/x/net/context"
)

type ContactRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewContactRepositoryBun(db *bun.DB) *ContactRepositoryBun {
	return &ContactRepositoryBun{db: db}
}

func (r *ContactRepositoryBun) CreateContact(ctx context.Context, c *personentity.Contact) error {
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

func (r *ContactRepositoryBun) UpdateContact(ctx context.Context, c *personentity.Contact) error {
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

	_, err := r.db.NewDelete().Model(&personentity.Contact{}).Where("id = ?", id).Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepositoryBun) GetContactById(ctx context.Context, id string) (*personentity.Contact, error) {
	contact := &personentity.Contact{}

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

func (r *ContactRepositoryBun) GetContactByDddAndNumber(ctx context.Context, ddd string, number string, contactType personentity.ContactType) (*personentity.Contact, error) {
	contact := &personentity.Contact{}

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

func (r *ContactRepositoryBun) FtSearchContacts(ctx context.Context, text string, contactType personentity.ContactType) (contacts []personentity.Contact, err error) {
	contacts = []personentity.Contact{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	err = r.db.NewSelect().Model(&contacts).Where("ts @@ websearch_to_tsquery('simple', ?) and type = ?", text, contactType).Scan(ctx)
	return
}
