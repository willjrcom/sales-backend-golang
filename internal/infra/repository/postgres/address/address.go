package addressrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type AddressRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewAddressRepositoryBun(db *bun.DB) *AddressRepositoryBun {
	return &AddressRepositoryBun{db: db}
}

func (r *AddressRepositoryBun) RegisterAddress(ctx context.Context, c *addressentity.Address) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(c).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepositoryBun) UpdateAddress(ctx context.Context, c *addressentity.Address) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepositoryBun) DeleteAddress(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewDelete().Model(&addressentity.Address{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *AddressRepositoryBun) GetAddressById(ctx context.Context, id string) (*addressentity.Address, error) {
	aAddress := &addressentity.Address{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(aAddress).Where("address.id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return aAddress, nil
}

func (r *AddressRepositoryBun) GetAllAddress(ctx context.Context) ([]addressentity.Address, error) {
	addresss := []addressentity.Address{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&addresss).Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return addresss, nil
}
