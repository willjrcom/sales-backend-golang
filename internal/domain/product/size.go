package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrSizeAlreadyExists = errors.New("size already exists")
)

type Size struct {
	entity.Entity
	SizeCommonAttributes
}

type SizeCommonAttributes struct {
	Name       string
	IsActive   *bool
	CategoryID uuid.UUID
	Products   []Product
}

func NewSize(sizeCommonAttributes SizeCommonAttributes) *Size {
	return &Size{
		Entity:               entity.NewEntity(),
		SizeCommonAttributes: sizeCommonAttributes,
	}
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
