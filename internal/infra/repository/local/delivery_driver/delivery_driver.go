package deliverydriverrepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type DeliveryDriverRepositoryLocal struct {}

func NewDeliveryDriverRepositoryLocal() model.DeliveryDriverRepository {
	return &DeliveryDriverRepositoryLocal{}
}

func (r *DeliveryDriverRepositoryLocal) CreateDeliveryDriver(ctx context.Context, p *model.DeliveryDriver) error {
	return nil
}

func (r *DeliveryDriverRepositoryLocal) UpdateDeliveryDriver(ctx context.Context, p *model.DeliveryDriver) error {
	return nil
}

func (r *DeliveryDriverRepositoryLocal) DeleteDeliveryDriver(ctx context.Context, id string) error {
	return nil
}

func (r *DeliveryDriverRepositoryLocal) GetDeliveryDriverById(ctx context.Context, id string) (*model.DeliveryDriver, error) {
	return nil, nil
}

func (r *DeliveryDriverRepositoryLocal) GetDeliveryDriverByEmployeeId(ctx context.Context, id string) (*model.DeliveryDriver, error) {
	return nil, nil
}

func (r *DeliveryDriverRepositoryLocal) GetAllDeliveryDrivers(ctx context.Context) ([]model.DeliveryDriver, error) {
	return nil, nil
}
