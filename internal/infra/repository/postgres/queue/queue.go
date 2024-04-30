package queuerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
)

type QueueRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewQueueRepositoryBun(db *bun.DB) *QueueRepositoryBun {
	return &QueueRepositoryBun{db: db}
}

func (r *QueueRepositoryBun) RegisterQueue(ctx context.Context, s *processentity.Queue) error {
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

func (r *QueueRepositoryBun) UpdateQueue(ctx context.Context, s *processentity.Queue) error {
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

	if _, err := r.db.NewDelete().Model(&processentity.Queue{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *QueueRepositoryBun) GetQueueById(ctx context.Context, id string) (*processentity.Queue, error) {
	queue := &processentity.Queue{}

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

func (r *QueueRepositoryBun) GetOpenedQueueByGroupItemId(ctx context.Context, id string) (*processentity.Queue, error) {
	queue := &processentity.Queue{}

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

func (r *QueueRepositoryBun) GetQueuesByGroupItemId(ctx context.Context, id string) ([]processentity.Queue, error) {
	queues := []processentity.Queue{}

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

func (r *QueueRepositoryBun) GetAllQueues(ctx context.Context) ([]processentity.Queue, error) {
	queuees := []processentity.Queue{}

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
