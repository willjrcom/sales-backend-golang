package employeeentity

import (
	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Employee struct {
	bun.BaseModel `bun:"table:employees"`
	personentity.Person
}
