package personentity

import (
	"time"

	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type Person struct {
	PersonCommonAttributes
}

type PersonCommonAttributes struct {
	ImagePath string
	Name      string
	Email     string
	Cpf       string
	Birthday  *time.Time
	IsActive  bool
	Contact   *Contact
	Address   *addressentity.Address
}

func NewPerson(personCommonAttributes *PersonCommonAttributes) *Person {
	return &Person{
		PersonCommonAttributes: *personCommonAttributes,
	}
}
