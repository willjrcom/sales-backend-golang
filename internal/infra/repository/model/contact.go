package model

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

var (
	ErrContactInvalid = errors.New("contact format invalid")
)

type Contact struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:contacts"`
	ContactCommonAttributes
	ObjectID uuid.UUID `bun:"object_id,type:uuid,notnull"`
}

type ContactCommonAttributes struct {
	Ddd    string `bun:"ddd,notnull"`
	Number string `bun:"number,notnull"`
	Type   string `bun:"type,notnull"`
}

func (c *Contact) FromDomain(contact *personentity.Contact) {
	*c = Contact{
		ContactCommonAttributes: ContactCommonAttributes{
			Ddd:    contact.Ddd,
			Number: contact.Number,
			Type:   string(contact.Type),
		},
	}
}

func (c *Contact) ToDomain() *personentity.Contact {
	if c == nil {
		return nil
	}
	return &personentity.Contact{
		ContactCommonAttributes: personentity.ContactCommonAttributes{
			Ddd:    c.Ddd,
			Number: c.Number,
			Type:   personentity.ContactType(c.Type),
		},
	}
}
