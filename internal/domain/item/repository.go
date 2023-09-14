package itementity

import "context"

type Repository interface {
	AddItemOrder(ctx context.Context, item *Item) error
	RemoveItemOrder(ctx context.Context, item *Item) error
	UpdateItemOrder(ctx context.Context, id string, item *Item) error
}
