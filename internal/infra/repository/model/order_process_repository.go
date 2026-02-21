package model

import (
	"context"
)

type OrderProcessRepository interface {
	CreateProcess(ctx context.Context, p *OrderProcess) error
	UpdateProcess(ctx context.Context, p *OrderProcess) error
	DeleteProcess(ctx context.Context, id string) error
	GetProcessById(ctx context.Context, id string) (*OrderProcess, error)
	GetAllProcessesFinishedByShiftID(ctx context.Context, shiftID string) ([]OrderProcess, error)
	GetProcessesByProcessRuleID(ctx context.Context, id string) ([]OrderProcess, error)
	GetProcessesByProductID(ctx context.Context, id string) ([]OrderProcess, error)
	GetProcessesByGroupItemID(ctx context.Context, id string) ([]OrderProcess, error)
	GetActiveProcessByGroupItemAndProcessRule(ctx context.Context, groupItemID, processRuleID string) (*OrderProcess, error)
}
