package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type PublicContact struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:public.contacts"`
	ContactCommonAttributes
	ObjectID uuid.UUID `bun:"object_id,type:uuid,notnull"`
}

func (c *PublicContact) FromDomain(contact *personentity.Contact) {
	if contact == nil {
		return
	}
	*c = PublicContact{
		Entity:   entitymodel.FromDomain(contact.Entity),
		ObjectID: contact.ObjectID,
		ContactCommonAttributes: ContactCommonAttributes{
			Ddd:    contact.Ddd,
			Number: contact.Number,
			Type:   string(contact.Type),
		},
	}
}

func (c *PublicContact) ToDomain() *personentity.Contact {
	if c == nil {
		return nil
	}
	return &personentity.Contact{
		Entity:   c.Entity.ToDomain(),
		ObjectID: c.ObjectID,
		ContactCommonAttributes: personentity.ContactCommonAttributes{
			Ddd:    c.Ddd,
			Number: c.Number,
			Type:   personentity.ContactType(c.Type),
		},
	}
}
