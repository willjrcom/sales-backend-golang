package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Category struct {
	entity.Entity
	bun.BaseModel       `bun:"table:categories"`
	Name                string     `bun:"name,notnull"`
	Sizes               []Size     `bun:"rel:has-many,join:id=category_id"`
	Quantities          []Quantity `bun:"rel:has-many,join:id=category_id"`
	Products            []Product  `bun:"rel:has-many,join:id=category_id"`
	Processes           []Process  `bun:"rel:has-many,join:id=category_id"`
	AditionalCategories []Category `bun:"m2m:category_to_aditional_categories,join:Category=Category"`
}

type CategoryToAditionalCategories struct {
	CategoryID  uuid.UUID `bun:"type:uuid,pk"`
	Category    *Category `bun:"rel:belongs-to,join:category_id=id"`
	AditionalID uuid.UUID `bun:"type:uuid,pk"`
	Aditional   *Category `bun:"rel:belongs-to,join:aditional_id=id"`
}
