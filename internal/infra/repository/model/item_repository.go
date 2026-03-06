package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ItemRepository interface {
	AddItemWithTx(ctx context.Context, tx *bun.Tx, item *Item) error
	AddAdditionalItemWithTx(ctx context.Context, tx *bun.Tx, id uuid.UUID, productID uuid.UUID, additionalItem *Item) error
	DeleteItem(ctx context.Context, id string) error
	DeleteItemWithTx(ctx context.Context, tx *bun.Tx, id string) error
	DeleteAdditionalItem(ctx context.Context, idAdditional uuid.UUID) error
	UpdateItem(ctx context.Context, item *Item) error
	UpdateItemWithTx(ctx context.Context, tx *bun.Tx, item *Item) error
	GetItemById(ctx context.Context, id string) (*Item, error)
	GetItemByIdWithTx(ctx context.Context, tx *bun.Tx, id string) (*Item, error)
	GetItemByAdditionalItemID(ctx context.Context, idAdditional uuid.UUID) (*Item, error)
}
