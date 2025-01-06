package model

import "context"

type QuantityRepository interface {
	CreateQuantity(ctx context.Context, Quantity *Quantity) error
	UpdateQuantity(ctx context.Context, Quantity *Quantity) error
	DeleteQuantity(ctx context.Context, id string) error
	GetQuantityById(ctx context.Context, id string) (*Quantity, error)
}
