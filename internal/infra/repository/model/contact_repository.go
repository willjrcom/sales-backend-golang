package model

import (
	"context"
)

type ContactRepository interface {
	CreateContact(ctx context.Context, c *Contact) (err error)
	UpdateContact(ctx context.Context, c *Contact) (err error)
	DeleteContact(ctx context.Context, id string) (err error)
	GetContactById(ctx context.Context, id string) (*Contact, error)
	GetContactByNumber(ctx context.Context, number string, contactType string) (*Contact, error)
	FtSearchContacts(ctx context.Context, key string, contactType string) (contacts []Contact, err error)
}
