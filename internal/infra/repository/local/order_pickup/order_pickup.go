package orderpickuprepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderPickupRepositoryLocal struct {}

func NewOrderPickupRepositoryLocal() model.OrderPickupRepository {
	return &OrderPickupRepositoryLocal{}
}

func (r *OrderPickupRepositoryLocal) CreateOrderPickup(ctx context.Context, pickup *model.OrderPickup) error {
	return nil
}

func (r *OrderPickupRepositoryLocal) UpdateOrderPickup(ctx context.Context, pickup *model.OrderPickup) error {
	return nil
}

func (r *OrderPickupRepositoryLocal) DeleteOrderPickup(ctx context.Context, id string) error {
	return nil
}

func (r *OrderPickupRepositoryLocal) GetPickupById(ctx context.Context, id string) (*model.OrderPickup, error) {
	return nil, nil
}

func (r *OrderPickupRepositoryLocal) GetAllPickups(ctx context.Context) ([]model.OrderPickup, error) {
	return nil, nil
}
