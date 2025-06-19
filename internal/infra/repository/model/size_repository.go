package model

import "context"

type SizeRepository interface {
	CreateSize(ctx context.Context, Size *Size) error
	UpdateSize(ctx context.Context, Size *Size) error
	DeleteSize(ctx context.Context, id string) error
	GetSizeById(ctx context.Context, id string) (*Size, error)
	GetSizeByIdWithProducts(ctx context.Context, id string) (*Size, error)
}
