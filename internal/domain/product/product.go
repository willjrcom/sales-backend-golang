package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrSizeIsInvalid    = errors.New("size is invalid")
)

type Product struct {
	entity.Entity
	bun.BaseModel `bun:"table:products"`
	Code          string    `bun:"code"`
	Name          string    `bun:"name,notnull"`
	Description   string    `bun:"description"`
	Price         float64   `bun:"price,notnull"`
	Cost          float64   `bun:"cost"`
	IsAvailable   bool      `bun:"is_available"`
	CategoryID    uuid.UUID `bun:"column:category_id,type:uuid,notnull"`
	Category      *Category `bun:"rel:belongs-to"`
	SizeID        uuid.UUID `bun:"column:size_id,type:uuid,notnull"`
	Size          *Size     `bun:"rel:belongs-to"`
}

func (p *Product) FindSizeInCategory() (bool, error) {
	if p.Category == nil {
		return false, ErrCategoryNotFound
	}

	for _, v := range p.Category.Sizes {
		if v.ID == p.SizeID {
			return true, nil
		}
	}

	return false, errors.New("size not found")
}
