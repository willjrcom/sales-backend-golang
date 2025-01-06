package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Employee struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:employees"`
	EmployeeCommonAttributes
}

type EmployeeCommonAttributes struct {
	UserID *uuid.UUID `bun:"column:user_id,type:uuid,notnull"`
	User   *User      `bun:"rel:belongs-to"`
}

func (e *Employee) FromDomain(employee *employeeentity.Employee) {
	*e = Employee{
		Entity: entitymodel.FromDomain(employee.Entity),
		EmployeeCommonAttributes: EmployeeCommonAttributes{
			UserID: employee.UserID,
			User:   &User{},
		},
	}

	e.User.FromDomain(employee.User)
}

func (e *Employee) ToDomain() *employeeentity.Employee {
	if e == nil {
		return nil
	}
	return &employeeentity.Employee{
		Entity: e.Entity.ToDomain(),
		UserID: e.UserID,
		User:   e.User.ToDomain(),
	}
}
