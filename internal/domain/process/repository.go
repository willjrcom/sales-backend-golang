package processentity

import (
	"context"
)

type ProcessRepository interface {
	RegisterProcess(ctx context.Context, p *Process) error
	UpdateProcess(ctx context.Context, p *Process) error
	DeleteProcess(ctx context.Context, id string) error
	GetProcessById(ctx context.Context, id string) (*Process, error)
	GetAllProcesses(ctx context.Context) ([]Process, error)
	GetProcessesByProductID(ctx context.Context, id string) ([]Process, error)
	GetProcessesByGroupItemID(ctx context.Context, id string) ([]Process, error)
}

type QueueRepository interface {
	RegisterQueue(ctx context.Context, p *Queue) error
	UpdateQueue(ctx context.Context, p *Queue) error
	DeleteQueue(ctx context.Context, id string) error
	GetQueueById(ctx context.Context, id string) (*Queue, error)
	GetQueueByPreviousProcessId(ctx context.Context, id string) (*Queue, error)
	GetAllQueuees(ctx context.Context) ([]Queue, error)
}
