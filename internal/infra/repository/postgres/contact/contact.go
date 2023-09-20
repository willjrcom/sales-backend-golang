package contactrepositorybun

import (
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
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

func (r *ContactRepositoryBun) GetContactsBy(ctx context.Context, c *personentity.Contact) ([]personentity.Contact, error) {
	contacts := []personentity.Contact{}

	r.mu.Lock()
	query := r.db.NewSelect().Model(&personentity.Contact{})

	if c.PersonID != uuid.Nil {
		query.Where("person_id = ?", c.PersonID)
	}
	if c.Ddd != "" {
		query.Where("ddd = ?", c.Ddd)
	}
	if c.Number != "" {
		query.Where("number LIKE ?", "%"+c.Number+"%")
	}

	err := query.Scan(ctx, &contacts)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *ContactRepositoryBun) GetAllContacts(ctx context.Context) ([]personentity.Contact, error) {
	Contacts := []personentity.Contact{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&Contacts).Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return Contacts, nil
}
