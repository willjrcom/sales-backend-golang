package placerepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type PlaceRepositoryBun struct {
	db *bun.DB
}

func NewPlaceRepositoryBun(db *bun.DB) model.PlaceRepository {
	return &PlaceRepositoryBun{db: db}
}

func (r *PlaceRepositoryBun) CreatePlace(ctx context.Context, s *model.Place) error {

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

func (r *PlaceRepositoryBun) UpdatePlace(ctx context.Context, s *model.Place) error {

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

func (r *PlaceRepositoryBun) DeletePlace(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Soft delete: set is_active to false
	isActive := false
	if _, err := tx.NewUpdate().
		Model(&model.Place{}).
		Set("is_active = ?", isActive).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *PlaceRepositoryBun) GetPlaceById(ctx context.Context, id string) (*model.Place, error) {
	place := &model.Place{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(place).Where("id = ?", id).Relation("Tables.Table").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return place, nil
}

func (r *PlaceRepositoryBun) GetAllPlaces(ctx context.Context, page, perPage int, isActive bool) ([]model.Place, int, error) {
	places := make([]model.Place, 0)

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	query := tx.NewSelect().
		Model(&places).
		Where("is_active = ?", isActive).
		Relation("Tables.Table").
		Limit(perPage).
		Offset(page * perPage)

	count, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return places, count, nil
}

func (r *PlaceRepositoryBun) AddTableToPlace(ctx context.Context, placeToTables *model.PlaceToTables) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.PlaceToTables{}).Where("table_id = ?", placeToTables.TableID).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewInsert().Model(placeToTables).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepositoryBun) GetTableToPlaceByTableID(ctx context.Context, tableID uuid.UUID) (*model.PlaceToTables, error) {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	placeToTable := &model.PlaceToTables{}

	if err := tx.NewSelect().Model(placeToTable).Where("table_id = ?", tableID).Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return placeToTable, nil
}

func (r *PlaceRepositoryBun) GetTableToPlaceByPlaceIDAndPosition(ctx context.Context, placeID uuid.UUID, column, row int) (*model.PlaceToTables, error) {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	placeToTable := &model.PlaceToTables{}

	if err := tx.NewSelect().Model(placeToTable).
		Where("place_id = ? AND \"column\" = ? AND row = ?", placeID.String(), column, row).Relation("Table").Relation("Place").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return placeToTable, nil
}

func (r *PlaceRepositoryBun) RemoveTableFromPlace(ctx context.Context, tableID uuid.UUID) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.PlaceToTables{}).Where("table_id = ?", tableID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
