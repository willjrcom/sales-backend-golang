package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrSizeAlreadyExists = errors.New("size already exists")
)

type Size struct {
	entity.Entity
	bun.BaseModel `bun:"table:sizes"`
	SizeCommonAttributes
}

type SizeCommonAttributes struct {
	Name       string    `bun:"name" json:"name"`
	Active     *bool     `bun:"active" json:"active"`
	CategoryID uuid.UUID `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Products   []Product `bun:"rel:has-many,join:id=size_id" json:"products,omitempty"`
}

type PatchSize struct {
	Name   *string `json:"name"`
	Active *bool   `json:"active"`
}

func ValidateDuplicateSizes(name string, sizes []Size) error {
	for _, size := range sizes {
		if size.Name == name {
			return ErrSizeAlreadyExists
		}
	}

	return nil
}

func ValidateUpdateSize(size *Size, sizes []Size) error {
	for _, s := range sizes {
		if s.Name == size.Name && s.ID != size.ID {
			return ErrSizeAlreadyExists
		}
	}

	return nil
}
