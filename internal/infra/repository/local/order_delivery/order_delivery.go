package orderdeliveryrepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderDeliveryRepositoryLocal struct {}

func NewOrderDeliveryRepositoryLocal() model.OrderDeliveryRepository {
	return &OrderDeliveryRepositoryLocal{}
}

func (r *OrderDeliveryRepositoryLocal) CreateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {
	return nil
}

func (r *OrderDeliveryRepositoryLocal) UpdateOrderDelivery(ctx context.Context, delivery *model.OrderDelivery) error {
	return nil
}

func (r *OrderDeliveryRepositoryLocal) DeleteOrderDelivery(ctx context.Context, id string) error {
	return nil
}

func (r *OrderDeliveryRepositoryLocal) GetDeliveryById(ctx context.Context, id string) (*model.OrderDelivery, error) {
	return nil, nil
}

func (r *OrderDeliveryRepositoryLocal) GetDeliveriesByIds(ctx context.Context, ids []string) ([]model.OrderDelivery, error) {
	return nil, nil
}

func (r *OrderDeliveryRepositoryLocal) GetAllDeliveries(ctx context.Context) ([]model.OrderDelivery, error) {
	return nil, nil
}
