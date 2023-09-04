package employeeentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Employee struct {
	entity.Entity
	personentity.Person
}
