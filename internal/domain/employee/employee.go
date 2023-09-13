package employeeentity

import (
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type Employee struct {
	personentity.Person
}
