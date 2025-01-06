package placerepositorybun

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type PlaceRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewPlaceRepositoryBun(db *bun.DB) *PlaceRepositoryBun {
	return &PlaceRepositoryBun{db: db}
}

func (r *PlaceRepositoryBun) CreatePlace(ctx context.Context, s *model.Place) error {
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

func (r *PlaceRepositoryBun) UpdatePlace(ctx context.Context, s *model.Place) error {
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

func (r *PlaceRepositoryBun) DeletePlace(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.Place{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepositoryBun) GetPlaceById(ctx context.Context, id string) (*model.Place, error) {
	place := &model.Place{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(place).Where("id = ?", id).Relation("Tables.Table").Scan(ctx); err != nil {
		return nil, err
	}

	return place, nil
}

func (r *PlaceRepositoryBun) GetAllPlaces(ctx context.Context) ([]model.Place, error) {
	places := make([]model.Place, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&places).Relation("Tables.Table").Scan(ctx); err != nil {
		return nil, err
	}

	return places, nil
}

func (r *PlaceRepositoryBun) AddTableToPlace(ctx context.Context, placeToTables *model.PlaceToTables) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.Begin()

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
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	placeToTable := &model.PlaceToTables{}

	if err := r.db.NewSelect().Model(placeToTable).Where("table_id = ?", tableID).Scan(ctx); err != nil {
		return nil, err
	}

	return placeToTable, nil
}

func (r *PlaceRepositoryBun) GetTableToPlaceByPlaceIDAndPosition(ctx context.Context, placeID uuid.UUID, column, row int) (*model.PlaceToTables, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	placeToTable := &model.PlaceToTables{}

	if err := r.db.NewSelect().Model(placeToTable).
		Where("place_id = ? AND \"column\" = ? AND row = ?", placeID.String(), column, row).Relation("Table").Relation("Place").Scan(ctx); err != nil {
		return nil, err
	}

	return placeToTable, nil
}

func (r *PlaceRepositoryBun) RemoveTableFromPlace(ctx context.Context, tableID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&model.PlaceToTables{}).Where("table_id = ?", tableID).Exec(ctx); err != nil {
		return err
	}

	return nil
}
