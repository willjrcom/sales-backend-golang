package personentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrContactInvalid = errors.New("contact format invalid")
)

type Contact struct {
	entity.Entity
	ContactCommonAttributes
	ObjectID uuid.UUID
}

type ContactCommonAttributes struct {
	Number string
	Type   ContactType
}

type ContactType string

const (
	ContactTypeClient   ContactType = "Client"
	ContactTypeEmployee ContactType = "Employee"
)

func GetAllOrderStatus() []ContactType {
	return []ContactType{
		ContactTypeClient,
		ContactTypeEmployee,
	}
}

func NewContact(contactCommonAttributes *ContactCommonAttributes) *Contact {
	return &Contact{
		Entity:                  entity.NewEntity(),
		ContactCommonAttributes: *contactCommonAttributes,
	}
}
