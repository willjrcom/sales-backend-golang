package model

import (
	"context"
)

type OrderProcessRepository interface {
	CreateProcess(ctx context.Context, p *OrderProcess) error
	UpdateProcess(ctx context.Context, p *OrderProcess) error
	DeleteProcess(ctx context.Context, id string) error
	GetProcessById(ctx context.Context, id string) (*OrderProcess, error)
	GetAllProcesses(ctx context.Context) ([]OrderProcess, error)
	GetProcessesByProcessRuleID(ctx context.Context, id string) ([]OrderProcess, error)
	GetProcessesByProductID(ctx context.Context, id string) ([]OrderProcess, error)
	GetProcessesByGroupItemID(ctx context.Context, id string) ([]OrderProcess, error)
}
