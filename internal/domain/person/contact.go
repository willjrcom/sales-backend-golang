package personentity

import (
	"errors"
	"regexp"
	"strings"

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
	Ddd    string
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

func ValidateAndExtractContact(text string) (ddd string, number string, err error) {
	// Use uma expressão regular para extrair o DDD e o número
	re := regexp.MustCompile(`\(?(\d{2})\)?\s*(\d{4,5}(?:-|\s)?\d{4})`)
	match := re.FindStringSubmatch(text)

	if len(match) < 3 {
		return "", "", ErrContactInvalid
	}

	ddd = match[1]
	number = match[2]
	number = strings.Replace(number, "-", "", -1)
	number = strings.Replace(number, " ", "", -1)
	err = nil
	return
}
