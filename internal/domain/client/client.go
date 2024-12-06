package cliententity

import (
	"time"

	"github.com/uptrace/bun"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Client struct {
	bun.BaseModel `bun:"table:clients"`
	personentity.Person
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
