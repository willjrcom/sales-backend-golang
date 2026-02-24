package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type AdvertisingRepository struct {
	db *bun.DB
}

func NewAdvertisingRepository(db *bun.DB) *AdvertisingRepository {
	return &AdvertisingRepository{db: db}
}

func (r *AdvertisingRepository) Create(ctx context.Context, advertising *model.Advertising) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(advertising).Exec(ctx); err != nil {
		return err
	}

	for _, cat := range advertising.Categories {
		categoryAdv := &model.CategoryToAdvertising{
			CategoryID:    cat.ID,
			AdvertisingID: advertising.ID,
		}
		if _, err := tx.NewInsert().Model(categoryAdv).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *AdvertisingRepository) Update(ctx context.Context, advertising *model.Advertising) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	res, err := tx.NewUpdate().Model(advertising).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("advertising not found")
	}

	if _, err := tx.NewDelete().Model((*model.CategoryToAdvertising)(nil)).Where("advertising_id = ?", advertising.ID).Exec(ctx); err != nil {
		return err
	}

	for _, cat := range advertising.Categories {
		categoryAdv := &model.CategoryToAdvertising{
			CategoryID:    cat.ID,
			AdvertisingID: advertising.ID,
		}
		if _, err := tx.NewInsert().Model(categoryAdv).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *AdvertisingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	res, err := tx.NewDelete().Model(&model.Advertising{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("advertising not found")
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *AdvertisingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Advertising, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	advertising := new(model.Advertising)
	err = tx.NewSelect().Model(advertising).Relation("Sponsor").Where("advertising.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return advertising, nil
}

func (r *AdvertisingRepository) GetAllAdvertisements(ctx context.Context) ([]model.Advertising, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	var advertisements []model.Advertising
	err = tx.NewSelect().Model(&advertisements).Relation("Sponsor").Order("title ASC").Scan(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return advertisements, nil
}
