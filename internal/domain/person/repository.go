package personentity

import (
	"context"
)

type ContactRepository interface {
	RegisterContact(ctx context.Context, c *Contact) error
	UpdateContact(ctx context.Context, c *Contact) error
	DeleteContact(ctx context.Context, id string) error
	GetContactById(ctx context.Context, id string) (*Contact, error)
	GetAllContacts(ctx context.Context) ([]Contact, error)
}
