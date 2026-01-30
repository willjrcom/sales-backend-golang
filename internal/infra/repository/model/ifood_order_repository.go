package model

import "context"

type IfoodOrderRepository interface {
	Create(ctx context.Context, o *IfoodOrder) error
	UpdateStatus(ctx context.Context, ifoodOrderID string, status string) error
	GetByIfoodOrderID(ctx context.Context, ifoodOrderID string) (*IfoodOrder, error)
}
