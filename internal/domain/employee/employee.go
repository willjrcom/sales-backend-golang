package employeeentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
)

type Employee struct {
	bun.BaseModel `bun:"table:employees"`
	personentity.Person
	UserID uuid.UUID        `bun:"column:user_id,type:uuid,notnull" json:"user_id"`
	User   *userentity.User `bun:"rel:belongs-to" json:"category,omitempty"`
}
