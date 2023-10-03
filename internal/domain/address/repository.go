package addressentity

import "context"

type Repository interface {
	RegisterAddress(ctx context.Context, address *Address) error
	UpdateAddress(ctx context.Context, address *Address) error
	RemoveAddress(ctx context.Context, id string) error
	GetAddressById(ctx context.Context, id string) (*Address, error)
}
