package personentity

import (
	"context"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type ContactRepository interface {
	RegisterContact(ctx context.Context, c *Contact) error
	UpdateContact(ctx context.Context, c *Contact) error
	DeleteContact(ctx context.Context, id string) error
	GetContactById(ctx context.Context, id string) (*Contact, error)
	GetContactsBy(ctx context.Context, c *Contact) ([]Contact, error)
	GetAllContacts(ctx context.Context) ([]Contact, error)
}

type AddressRepository interface {
	RegisterAddress(ctx context.Context, a *addressentity.Address) error
	UpdateAddress(ctx context.Context, a *addressentity.Address) error
	DeleteAddress(ctx context.Context, id string) error
}
