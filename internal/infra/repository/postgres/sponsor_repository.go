package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type SponsorRepository struct {
	db *bun.DB
}

func NewSponsorRepository(db *bun.DB) *SponsorRepository {
	return &SponsorRepository{db: db}
}

func (r *SponsorRepository) Create(ctx context.Context, sponsor *model.Sponsor) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(sponsor).Exec(ctx); err != nil {
		return err
	}

	if sponsor.Address != nil {
		sponsor.Address.ObjectID = sponsor.ID
		if _, err := tx.NewInsert().Model(sponsor.Address).Exec(ctx); err != nil {
			return err
		}
	}

	for _, cat := range sponsor.Categories {
		categorySponsor := &model.CategoryToSponsor{
			CategoryID: cat.ID,
			SponsorID:  sponsor.ID,
		}
		if _, err := tx.NewInsert().Model(categorySponsor).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *SponsorRepository) Update(ctx context.Context, sponsor *model.Sponsor) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	res, err := tx.NewUpdate().Model(sponsor).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("sponsor not found")
	}

	if sponsor.Address != nil {
		sponsor.Address.ObjectID = sponsor.ID
		if _, err := tx.NewUpdate().Model(sponsor.Address).WherePK().Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *SponsorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	res, err := tx.NewDelete().Model(&model.Sponsor{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("sponsor not found")
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *SponsorRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Sponsor, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	sponsor := new(model.Sponsor)
	err = tx.NewSelect().Model(sponsor).Relation("Address").Where("sponsor.id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sponsor, nil
}

func (r *SponsorRepository) GetAllSponsors(ctx context.Context) ([]model.Sponsor, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	var sponsors []model.Sponsor
	err = tx.NewSelect().Model(&sponsors).Relation("Address").Order("name ASC").Scan(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return sponsors, nil
}
