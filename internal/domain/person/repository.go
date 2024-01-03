package personentity

import (
	"context"
)

type ContactRepository interface {
	RegisterContact(ctx context.Context, c *Contact) (err error)
	UpdateContact(ctx context.Context, c *Contact) (err error)
	DeleteContact(ctx context.Context, id string) (err error)
	GetContactById(ctx context.Context, id string) (*Contact, error)
	GetAllContacts(ctx context.Context) (contacts []Contact, err error)
	FtSearchContacts(ctx context.Context, keys string) (contacts []Contact, err error)
}
