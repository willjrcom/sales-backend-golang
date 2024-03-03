package personentity

import (
	"context"
)

type ContactRepository interface {
	RegisterContact(ctx context.Context, c *Contact) (err error)
	UpdateContact(ctx context.Context, c *Contact) (err error)
	DeleteContact(ctx context.Context, id string) (err error)
	GetContactById(ctx context.Context, id string) (*Contact, error)
	GetContactByDddAndNumber(ctx context.Context, ddd string, number string, contactType ContactType) (*Contact, error)
	FtSearchContacts(ctx context.Context, key string, contactType ContactType) (contacts []Contact, err error)
}
