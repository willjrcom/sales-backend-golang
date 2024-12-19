package orderprocessentity

import (
	"context"
)

type ProcessRepository interface {
	CreateProcess(ctx context.Context, p *OrderProcess) error
	UpdateProcess(ctx context.Context, p *OrderProcess) error
	DeleteProcess(ctx context.Context, id string) error
	GetProcessById(ctx context.Context, id string) (*OrderProcess, error)
	GetAllProcesses(ctx context.Context) ([]OrderProcess, error)
	GetProcessesByProcessRuleID(ctx context.Context, id string) ([]OrderProcess, error)
	GetProcessesByProductID(ctx context.Context, id string) ([]OrderProcess, error)
	GetProcessesByGroupItemID(ctx context.Context, id string) ([]OrderProcess, error)
}

type QueueRepository interface {
	CreateQueue(ctx context.Context, p *OrderQueue) error
	UpdateQueue(ctx context.Context, p *OrderQueue) error
	DeleteQueue(ctx context.Context, id string) error
	GetQueueById(ctx context.Context, id string) (*OrderQueue, error)
	GetOpenedQueueByGroupItemId(ctx context.Context, id string) (*OrderQueue, error)
	GetQueuesByGroupItemId(ctx context.Context, id string) ([]OrderQueue, error)
	GetAllQueues(ctx context.Context) ([]OrderQueue, error)
}
