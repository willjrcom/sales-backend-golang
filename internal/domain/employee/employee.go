package employeeentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Employee struct {
	bun.BaseModel `bun:"table:employees"`
	personentity.Person
	UserID uuid.UUID           `bun:"column:user_id,type:uuid,notnull" json:"user_id"`
	User   *companyentity.User `bun:"rel:belongs-to" json:"category,omitempty"`
}
