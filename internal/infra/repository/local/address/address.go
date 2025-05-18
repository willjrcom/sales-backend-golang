package addressrepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type AddressRepositoryLocal struct {}

func NewAddressRepositoryLocal() model.AddressRepository {
	return &AddressRepositoryLocal{}
}

func (r *AddressRepositoryLocal) CreateAddress(ctx context.Context, address *model.Address) error {
	return nil
}

func (r *AddressRepositoryLocal) UpdateAddress(ctx context.Context, address *model.Address) error {
	return nil
}

func (r *AddressRepositoryLocal) DeleteAddress(ctx context.Context, id string) error {
	return nil
}

func (r *AddressRepositoryLocal) GetAddressById(ctx context.Context, id string) (*model.Address, error) {
	return nil, nil
}
