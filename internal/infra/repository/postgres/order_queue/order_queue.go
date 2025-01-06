package orderqueuerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type QueueRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderQueueRepositoryBun(db *bun.DB) *QueueRepositoryBun {
	return &QueueRepositoryBun{db: db}
}

func (r *QueueRepositoryBun) CreateQueue(ctx context.Context, s *model.OrderQueue) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *QueueRepositoryBun) UpdateQueue(ctx context.Context, s *model.OrderQueue) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *QueueRepositoryBun) DeleteQueue(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.OrderQueue{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *QueueRepositoryBun) GetQueueById(ctx context.Context, id string) (*model.OrderQueue, error) {
	queue := &model.OrderQueue{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(queue).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return queue, nil
}

func (r *QueueRepositoryBun) GetOpenedQueueByGroupItemId(ctx context.Context, id string) (*model.OrderQueue, error) {
	queue := &model.OrderQueue{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(queue).Where("group_item_id = ? AND left_at IS NULL", id).Scan(ctx); err != nil {
		return nil, err
	}

	return queue, nil
}

func (r *QueueRepositoryBun) GetQueuesByGroupItemId(ctx context.Context, id string) ([]model.OrderQueue, error) {
	queues := []model.OrderQueue{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&queues).Where("group_item_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return queues, nil
}

func (r *QueueRepositoryBun) GetAllQueues(ctx context.Context) ([]model.OrderQueue, error) {
	queuees := []model.OrderQueue{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&queuees).Scan(ctx); err != nil {
		return nil, err
	}

	return queuees, nil
}
