package model

import "context"

type DeliveryDriverRepository interface {
	CreateDeliveryDriver(ctx context.Context, DeliveryDriver *DeliveryDriver) error
	UpdateDeliveryDriver(ctx context.Context, DeliveryDriver *DeliveryDriver) error
	DeleteDeliveryDriver(ctx context.Context, id string) error
	GetDeliveryDriverById(ctx context.Context, id string) (*DeliveryDriver, error)
	GetDeliveryDriverByEmployeeId(ctx context.Context, id string) (*DeliveryDriver, error)
	GetAllDeliveryDrivers(ctx context.Context, isActive ...bool) ([]DeliveryDriver, error)
}
