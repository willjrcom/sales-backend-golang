package model

import "context"

type AddressRepository interface {
	CreateAddress(ctx context.Context, address *Address) error
	UpdateAddress(ctx context.Context, address *Address) error
	DeleteAddress(ctx context.Context, id string) error
	GetAddressById(ctx context.Context, id string) (*Address, error)
}
