package contactrepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ContactRepositoryLocal struct {}

func NewContactRepositoryLocal() model.ContactRepository {
	return &ContactRepositoryLocal{}
}

func (r *ContactRepositoryLocal) CreateContact(ctx context.Context, c *model.Contact) error {
	return nil
}

func (r *ContactRepositoryLocal) UpdateContact(ctx context.Context, c *model.Contact) error {
	return nil
}

func (r *ContactRepositoryLocal) DeleteContact(ctx context.Context, id string) error {
	return nil
}

func (r *ContactRepositoryLocal) GetContactById(ctx context.Context, id string) (*model.Contact, error) {
	return nil, nil
}

func (r *ContactRepositoryLocal) GetContactByDddAndNumber(ctx context.Context, ddd string, number string, contactType string) (*model.Contact, error) {
	return nil, nil
}

func (r *ContactRepositoryLocal) FtSearchContacts(ctx context.Context, key string, contactType string) ([]model.Contact, error) {
	return nil, nil
}
