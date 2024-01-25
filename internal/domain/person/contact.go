package personentity

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrContactInvalid = errors.New("contact format invalid")
)

type Contact struct {
	entity.Entity
	bun.BaseModel `bun:"table:contacts"`
	ContactCommonAttributes
}

type ContactCommonAttributes struct {
	ObjectID uuid.UUID   `bun:"object_id,type:uuid,notnull" json:"object_id"`
	Ddd      string      `bun:"ddd,notnull" json:"ddd"`
	Number   string      `bun:"number,notnull" json:"number"`
	Type     ContactType `bun:"type,notnull" json:"type"`
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
