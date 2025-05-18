package queuerepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type QueueRepositoryLocal struct {}

func NewQueueRepositoryLocal() model.QueueRepository {
	return &QueueRepositoryLocal{}
}

func (r *QueueRepositoryLocal) CreateQueue(ctx context.Context, p *model.OrderQueue) error {
	return nil
}

func (r *QueueRepositoryLocal) UpdateQueue(ctx context.Context, p *model.OrderQueue) error {
	return nil
}

func (r *QueueRepositoryLocal) DeleteQueue(ctx context.Context, id string) error {
	return nil
}

func (r *QueueRepositoryLocal) GetQueueById(ctx context.Context, id string) (*model.OrderQueue, error) {
	return nil, nil
}

func (r *QueueRepositoryLocal) GetOpenedQueueByGroupItemId(ctx context.Context, id string) (*model.OrderQueue, error) {
	return nil, nil
}

func (r *QueueRepositoryLocal) GetQueuesByGroupItemId(ctx context.Context, id string) ([]model.OrderQueue, error) {
	return nil, nil
}

func (r *QueueRepositoryLocal) GetAllQueues(ctx context.Context) ([]model.OrderQueue, error) {
	return nil, nil
}
