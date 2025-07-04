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
	Payments    []PaymentEmployee `bun:"rel:has-many,join:id=employee_id"`
	Permissions map[string]string `bun:"type:jsonb"`
}

type EmployeeCommonAttributes struct {
	UserID uuid.UUID `bun:"column:user_id,type:uuid,notnull"`
	User   *User     `bun:"rel:belongs-to,join:user_id=id"`
}

func (e *Employee) FromDomain(employee *employeeentity.Employee) {
	if employee == nil {
		return
	}
	*e = Employee{
		Entity: entitymodel.FromDomain(employee.Entity),
		EmployeeCommonAttributes: EmployeeCommonAttributes{
			UserID: employee.UserID,
			User:   &User{},
		},
		Payments:    []PaymentEmployee{},
		Permissions: make(map[string]string),
	}

	e.User.FromDomain(employee.User)
	// map payments from domain
	for _, pay := range employee.Payments {
		p := PaymentEmployee{}
		p.FromDomain(&pay)
		e.Payments = append(e.Payments, p)
	}
	// map permissions from domain
	for k, v := range employee.Permissions {
		e.Permissions[string(k)] = v
	}
}

func (e *Employee) ToDomain() *employeeentity.Employee {
	if e == nil {
		return nil
	}
	dom := &employeeentity.Employee{
		Entity:      e.Entity.ToDomain(),
		UserID:      e.UserID,
		User:        e.User.ToDomain(),
		Payments:    make([]employeeentity.PaymentEmployee, 0, len(e.Payments)),
		Permissions: make(employeeentity.Permissions),
	}
	// map payments to domain
	for _, p := range e.Payments {
		dom.Payments = append(dom.Payments, *p.ToDomain())
	}
	// map permissions to domain
	for k, v := range e.Permissions {
		dom.Permissions[employeeentity.PermissionKey(k)] = v
	}
	return dom
}
