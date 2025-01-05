package model

import "context"

type SchemaRepository interface {
	NewSchema(ctx context.Context) error
}
