package itementity

import (
	"context"
)

type ItemRepository interface {
	AddItem(ctx context.Context, item *Item) error
	DeleteItem(ctx context.Context, id string) error
	UpdateItem(ctx context.Context, item *Item) error
}
