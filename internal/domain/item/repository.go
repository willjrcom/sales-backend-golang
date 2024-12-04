package itementity

import (
	"context"

	"github.com/google/uuid"
)

type ItemRepository interface {
	AddItem(ctx context.Context, item *Item) error
	AddAdditionalItem(ctx context.Context, id uuid.UUID, additionalItem *Item) error
	DeleteItem(ctx context.Context, id string) error
	DeleteAdditionalItem(ctx context.Context, idAdditional uuid.UUID) error
	UpdateItem(ctx context.Context, item *Item) error
	GetItemById(ctx context.Context, id string) (*Item, error)
}
