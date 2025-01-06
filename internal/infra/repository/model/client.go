package model

import (
	"github.com/uptrace/bun"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Client struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:clients"`
	Person
}
