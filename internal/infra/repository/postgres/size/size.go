package sizerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type SizeRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewSizeRepositoryBun(db *bun.DB) model.SizeRepository {
	return &SizeRepositoryBun{db: db}
}

func (r *SizeRepositoryBun) CreateSize(ctx context.Context, s *model.Size) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *SizeRepositoryBun) UpdateSize(ctx context.Context, s *model.Size) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *SizeRepositoryBun) DeleteSize(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Size{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *SizeRepositoryBun) GetSizeById(ctx context.Context, id string) (*model.Size, error) {
	size := &model.Size{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(size).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return size, nil
}

func (r *SizeRepositoryBun) GetSizeByIdWithProducts(ctx context.Context, id string) (*model.Size, error) {
	size := &model.Size{}

	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(size).Where("id = ?", id).Relation("Products").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return size, nil
}
