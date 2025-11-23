package orderqueuerepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type QueueRepositoryBun struct {
	db *bun.DB
}

func NewOrderQueueRepositoryBun(db *bun.DB) model.QueueRepository {
	return &QueueRepositoryBun{db: db}
}

func (r *QueueRepositoryBun) CreateQueue(ctx context.Context, s *model.OrderQueue) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *QueueRepositoryBun) UpdateQueue(ctx context.Context, s *model.OrderQueue) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *QueueRepositoryBun) DeleteQueue(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.OrderQueue{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *QueueRepositoryBun) GetQueueById(ctx context.Context, id string) (*model.OrderQueue, error) {
	queue := &model.OrderQueue{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(queue).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return queue, nil
}

func (r *QueueRepositoryBun) GetOpenedQueueByGroupItemId(ctx context.Context, id string) (*model.OrderQueue, error) {
	queue := &model.OrderQueue{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(queue).Where("group_item_id = ? AND left_at IS NULL", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return queue, nil
}

func (r *QueueRepositoryBun) GetQueuesByGroupItemId(ctx context.Context, id string) ([]model.OrderQueue, error) {
	queues := []model.OrderQueue{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&queues).Where("group_item_id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return queues, nil
}

func (r *QueueRepositoryBun) GetAllQueues(ctx context.Context) ([]model.OrderQueue, error) {
	queuees := []model.OrderQueue{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(&queuees).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return queuees, nil
}
