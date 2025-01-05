package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProcessRule struct {
	entity.Entity
	bun.BaseModel `bun:"table:process_rules,alias:pr"`
	ProcessRuleCommonAttributes
}

type ProcessRuleCommonAttributes struct {
	Name              string        `bun:"name,notnull"`
	Order             int8          `bun:"order,notnull"`
	Description       string        `bun:"description"`
	ImagePath         *string       `bun:"image_path"`
	IdealTime         time.Duration `bun:"ideal_time,notnull"`
	ExperimentalError time.Duration `bun:"experimental_error,notnull"`
	CategoryID        uuid.UUID     `bun:"column:category_id,type:uuid,notnull"`
	DeletedAt         time.Time     `bun:",soft_delete,nullzero"`
}
