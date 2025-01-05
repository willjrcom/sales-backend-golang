package employeeentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Employee struct {
	entity.Entity
	bun.BaseModel `bun:"table:employees"`
	UserID        *uuid.UUID          `bun:"column:user_id,type:uuid" json:"user_id,omitempty"`
	User          *companyentity.User `bun:"rel:belongs-to" json:"user,omitempty"`
}
