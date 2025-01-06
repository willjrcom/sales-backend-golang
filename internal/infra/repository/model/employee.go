package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
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
