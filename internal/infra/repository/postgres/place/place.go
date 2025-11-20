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

func (r *PlaceRepositoryBun) UpdatePlace(ctx context.Context, s *model.Place) error {

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

func (r *PlaceRepositoryBun) DeletePlace(ctx context.Context, id string) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Place{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *PlaceRepositoryBun) GetPlaceById(ctx context.Context, id string) (*model.Place, error) {
	place := &model.Place{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(place).Where("id = ?", id).Relation("Tables.Table").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return place, nil
}

func (r *PlaceRepositoryBun) GetAllPlaces(ctx context.Context) ([]model.Place, error) {
	places := make([]model.Place, 0)

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(&places).Relation("Tables.Table").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return places, nil
}

func (r *PlaceRepositoryBun) AddTableToPlace(ctx context.Context, placeToTables *model.PlaceToTables) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.PlaceToTables{}).Where("table_id = ?", placeToTables.TableID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewInsert().Model(placeToTables).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepositoryBun) GetTableToPlaceByTableID(ctx context.Context, tableID uuid.UUID) (*model.PlaceToTables, error) {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

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

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

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

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.PlaceToTables{}).Where("table_id = ?", tableID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
