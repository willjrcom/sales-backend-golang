package contactrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	"golang.org/x/net/context"
)

type ContactRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewContactRepositoryBun(ctx context.Context, db *bun.DB) *ContactRepositoryBun {
	return &ContactRepositoryBun{db: db}
}

func (r *ContactRepositoryBun) RegisterContact(ctx context.Context, c *personentity.Contact) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(c).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepositoryBun) UpdateContact(ctx context.Context, c *personentity.Contact) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepositoryBun) DeleteContact(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&personentity.Contact{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepositoryBun) GetContactById(ctx context.Context, id string) (*personentity.Contact, error) {
	contact := &personentity.Contact{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(contact).Where("contact.id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *ContactRepositoryBun) FtSearchContacts(ctx context.Context, id string) (contacts []personentity.Contact, err error) {
	contacts = []personentity.Contact{}
	err = r.db.NewSelect().Model(&contacts).Where("ts @@ websearch_to_tsquery('simple', ?)", id).Scan(ctx)
	return contacts, err
}
