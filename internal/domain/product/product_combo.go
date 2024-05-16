package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductToCombo struct {
	bun.BaseModel  `bun:"table:product_category_to_combo"`
	ProductID      uuid.UUID `bun:"type:uuid,pk"`
	Product        *Product  `bun:"rel:belongs-to,join:product_id=id"`
	ComboProductID uuid.UUID `bun:"type:uuid,pk"`
	ComboProduct   *Product  `bun:"rel:belongs-to,join:combo_product_id=id"`
}
