package orderprocessrepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderProcessRepositoryLocal struct {}

func NewOrderProcessRepositoryLocal() model.OrderProcessRepository {
	return &OrderProcessRepositoryLocal{}
}

func (r *OrderProcessRepositoryLocal) CreateProcess(ctx context.Context, p *model.OrderProcess) error {
	return nil
}

func (r *OrderProcessRepositoryLocal) UpdateProcess(ctx context.Context, p *model.OrderProcess) error {
	return nil
}

func (r *OrderProcessRepositoryLocal) DeleteProcess(ctx context.Context, id string) error {
	return nil
}

func (r *OrderProcessRepositoryLocal) GetProcessById(ctx context.Context, id string) (*model.OrderProcess, error) {
	return nil, nil
}

func (r *OrderProcessRepositoryLocal) GetAllProcesses(ctx context.Context) ([]model.OrderProcess, error) {
	return nil, nil
}

func (r *OrderProcessRepositoryLocal) GetProcessesByProcessRuleID(ctx context.Context, id string) ([]model.OrderProcess, error) {
	return nil, nil
}

func (r *OrderProcessRepositoryLocal) GetProcessesByProductID(ctx context.Context, id string) ([]model.OrderProcess, error) {
	return nil, nil
}

func (r *OrderProcessRepositoryLocal) GetProcessesByGroupItemID(ctx context.Context, id string) ([]model.OrderProcess, error) {
	return nil, nil
}
