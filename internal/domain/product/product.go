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
	Code          string           `bun:"code"`
	Name          string           `bun:"name"`
	Description   string           `bun:"description"`
	Size          string           `bun:"size"`
	Price         float64          `bun:"price"`
	Cost          float64          `bun:"cost"`
	CategoryID    uuid.UUID        `bun:"column:category_id,type:uuid,notnull"`
	Category      *CategoryProduct `bun:"rel:belongs-to"`
	IsAvailable   bool             `bun:"is_available"`
}

func (p *Product) FindSizeInCategory() (bool, error) {
	if p.Category == nil {
		return false, ErrCategoryNotFound
	}

	for _, v := range p.Category.Sizes {
		if v == p.Size {
			return true, nil
		}
	}

	return false, nil
}
