package processentity

import (
	"context"
)

type ProcessRepository interface {
	RegisterProcess(ctx context.Context, p *Process) error
	UpdateProcess(ctx context.Context, p *Process) error
	DeleteProcess(ctx context.Context, id string) error
	GetProcessById(ctx context.Context, id string) (*Process, error)
	GetAllProcesss(ctx context.Context) ([]Process, error)
}
