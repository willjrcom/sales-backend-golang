package schemaentity

import "context"

type Repository interface {
	NewSchema(ctx context.Context, id string)
}
