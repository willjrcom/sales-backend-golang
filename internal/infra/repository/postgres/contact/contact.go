package contactrepositorybun

import (
	"log"
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
	setupFtSearch(ctx, db)
	return &ContactRepositoryBun{db: db}
}

func setupFtSearch(ctx context.Context, db *bun.DB) {
	column := `
		ALTER TABLE contacts ADD COLUMN IF NOT EXISTS ts tsvector
		GENERATED ALWAYS AS (
			setweight(to_tsvector('simple', coalesce(ddd::text, '') || coalesce(number::text, '')), 'A') ||
			setweight(to_tsvector('simple', coalesce(ddd::text, '')), 'B') ||
			setweight(to_tsvector('simple', coalesce(number::text, '')), 'B')
		) STORED;
	`

	_, err := db.ExecContext(ctx, column)
	if err != nil {
		log.Fatalln("Failed to create tsvector column for contacts", err)
	}

	index := "CREATE INDEX IF NOT EXISTS contacts_ts_idx ON contacts USING GIN(ts);"
	_, err = db.ExecContext(ctx, index)
	if err != nil {
		log.Fatalln("Failed to create index for contacts tsvector column")
	}

	log.Println("Created tsvector column and index for contacts table")
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
