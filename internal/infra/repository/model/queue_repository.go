package model

import "context"

type QueueRepository interface {
	CreateQueue(ctx context.Context, p *OrderQueue) error
	UpdateQueue(ctx context.Context, p *OrderQueue) error
	DeleteQueue(ctx context.Context, id string) error
	GetQueueById(ctx context.Context, id string) (*OrderQueue, error)
	GetOpenedQueueByGroupItemId(ctx context.Context, id string) (*OrderQueue, error)
	GetQueuesByGroupItemId(ctx context.Context, id string) ([]OrderQueue, error)
	GetAllQueues(ctx context.Context) ([]OrderQueue, error)
}
